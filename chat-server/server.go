package main

import (
	"os"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"crypto/rsa"
	"crypto/x509"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

var pubKey *rsa.PublicKey

func readJWTKey(jwtKeyFilename string) *rsa.PublicKey {
	raw, err := os.ReadFile(jwtKeyFilename)
	if err != nil {
		panic("failed to read public key file" + err.Error())
	}
	pub, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	return pub.(*rsa.PublicKey)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
	})

	return token, err
}

func verifyUser(r *http.Request) (jwt.MapClaims, bool) {
	jwtid, err := r.Cookie("JWTID")

	if err != nil {
		log.Println("Cookie \"JWTID\" not set")
		return nil, false
	}

	token, err := parseToken(jwtid.Value)
	if err != nil{
		log.Println(err)
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && !token.Valid {
		log.Println("Cookie \"JWTID\" is not valid")
		return nil, false
	}

	return claims, true
}

type Message struct {
    Messagetime time.Time
    Messagetext string
    Userid int
    Channelid int
}

type User struct {
	Userid int
    Username string
    Userhash string
    Usersalt string
    Useremail string
}

type ChatMessage struct {
	Username string
	Text string
	Timestamp string
}

type Payload struct {
	Op string
	Value string
}

type Client struct{
	Userid int
	Conn *websocket.Conn
}

type Channel struct {
	Clients map[int]Client
	ChannelId int
}

func userIdFromUsername(db *gorm.DB, username string) int{
	var user User
	db.Where("username = ?",username).Find(&user)
	return user.Userid
}

func createHandler(channels *map[int]Channel, messageQ chan<- *Message, upgrader websocket.Upgrader, db *gorm.DB) func(http.ResponseWriter,*http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		claims, auth := verifyUser(r)
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		upgrader.CheckOrigin = func(r *http.Request) bool { 
			return true 
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		log.Println("New connection from",conn.RemoteAddr().String())
		
		var client Client = Client{
			Userid: userIdFromUsername(db, claims["sub"].(string)),
		}
		var channel Channel
		var inChannel bool = false

		for {
			var p Payload
			err := conn.ReadJSON(&p)
			if err != nil {
				log.Println(err)
				return
			}
			
			switch p.Op{
			case "join":
				channelId, err := strconv.Atoi(p.Value)
				if err != nil {
					log.Println(err)
					return
				}

				val, channelExists := (*channels)[channelId]
				channel = val
				client.Conn = conn
				if channelExists {
					channel.Clients[client.Userid] = client
				}else{
					channel = Channel{
						Clients: map[int]Client{client.Userid: client},
						ChannelId: channelId,
					}
					(*channels)[channelId] = channel
				}

				inChannel = true
			case "get":
				if inChannel {
					msgNum, err := strconv.Atoi(p.Value)
					if err != nil {
						log.Println(err)
						break
					}
					var msgs []Message
					db.Where("channelid = ?",channel.ChannelId).Order("messagetime DESC").Limit(msgNum).Offset(0).Find(&msgs)
					
					var chatMsgs [] ChatMessage
					for i := range msgs {
						var user User
						db.First(&user, msgs[i].Userid)
						chatMsgs = append(chatMsgs, ChatMessage{
							Username: user.Username,
							Text: msgs[i].Messagetext,
							Timestamp: msgs[i].Messagetime.Format("2006-01-02 15:04:05"),
						})
					}
					b, err := json.Marshal(chatMsgs)
					if err != nil {
						log.Println("Failed to convert msg struct to json")
						break
					}

					client.Conn.WriteMessage(websocket.TextMessage, b)
				}
			case "send":
				if inChannel {
					msg := Message{
						time.Now(),
						string(p.Value),
						client.Userid,
						channel.ChannelId,
					}
					messageQ <- &msg
				}else {
					log.Printf("Client %d is not in Channel %d",client.Userid,channel.ChannelId)
					return
				}
			case "leave":
				channelId, err := strconv.Atoi(p.Value)
				if err != nil {
					log.Println(err)
					return
				}

				channel, exists := (*channels)[channelId]
				if exists {
					delete(channel.Clients, client.Userid)
				}

				log.Printf("Client %d left channel %d",client.Userid,channelId)

				return
			}
		}
	}
}

func channelLoop(db *gorm.DB, connList *map[int]Channel, messageQ <-chan *Message){
	for {
		msg := <- messageQ
		// persist
		go db.Create(&msg)
		// send to chat channel
		channel := (*connList)[msg.Channelid]
		clients := channel.Clients
		for clientId := range clients {
			conn := clients[clientId].Conn
			
			var user User
			db.First(&user, msg.Userid)
			b, err := json.Marshal(ChatMessage{
				Username: user.Username,
				Text: msg.Messagetext,
				Timestamp: msg.Messagetime.Format("2006-01-02 15:04:05"),
			})
			if err != nil {
				log.Println("Failed to convert msg struct to json")
				break
			}

			err = (*conn).WriteMessage(websocket.TextMessage, b)
			if err != nil{
				delete(clients,clientId)
				log.Printf("Removing Client %d from Channel %d: %s",clientId,channel.ChannelId,err)
			}
		}
	}
}

func getEnvVar(name string, dflt string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return dflt
}

func main() {
	jwtKeyFilename := getEnvVar("JWTKEY_FILENAME", "public_key.der")
	serverHost := getEnvVar("SERVER_HOST", "localhost")
	serverPort := getEnvVar("SERVER_PORT", "9000")
	dbHost := getEnvVar("POSTGRES_HOST", "localhost")
	dbPort := getEnvVar("POSTGRES_PORT", "5432")
	dbName := getEnvVar("POSTGRES_DB", "gochat")
	dbUser := getEnvVar("POSTGRES_USERNAME", "gochat")
	dbPassword := getEnvVar("POSTGRES_PASSWORD", "gochat")

	pubKey = readJWTKey(jwtKeyFilename);
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	var upgrader = websocket.Upgrader{ReadBufferSize:  1024, WriteBufferSize: 1024}
	var connList map[int]Channel = make(map[int]Channel)
	var messageQ chan *Message = make(chan *Message)

	go channelLoop(db, &connList, messageQ)

	handler := createHandler(&connList, messageQ, upgrader, db)
	http.HandleFunc("/ws", handler)
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%s",serverHost,serverPort),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Println("Listening...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}