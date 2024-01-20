package server

import (
	"log"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"crypto/rsa"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"

	"chatserver/internal/utils"
	"chatserver/internal/model"
)

type Payload struct {
	Op string
	Value string
}

type Message struct {
    Messagetime time.Time
    Messagetext string
    Userid int
    Channelid int
}

type Client struct{
	Userid int
	Conn *websocket.Conn
}

type Channel struct {
	Clients map[int]Client
	ChannelId int
}

type Chatserver struct {
	serverHost string
	serverPort string
	db *gorm.DB
	messageQ chan *Message
	channels map[int]Channel
	pubKey *rsa.PublicKey
}

func (chatserver *Chatserver) parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return chatserver.pubKey, nil
	})

	return token, err
}

func (chatserver *Chatserver) verifyUser(r *http.Request) (jwt.MapClaims, bool) {
	jwtid, err := r.Cookie("JWTID")

	if err != nil {
		log.Println("Cookie \"JWTID\" not set")
		return nil, false
	}

	token, err := chatserver.parseToken(jwtid.Value)
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

func (chatserver *Chatserver) userIdFromUsername(username string) int {
	var user model.User
	chatserver.db.Where("username = ?",username).Find(&user)
	return user.Userid
}

func (chatserver *Chatserver) createHandler() func(http.ResponseWriter,*http.Request) {
	var upgrader = websocket.Upgrader{ReadBufferSize:  1024, WriteBufferSize: 1024}
	
	return func (w http.ResponseWriter, r *http.Request) {
		claims, auth := chatserver.verifyUser(r)
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
			Userid: chatserver.userIdFromUsername(claims["sub"].(string)),
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

				val, channelExists := chatserver.channels[channelId]
				channel = val
				client.Conn = conn
				if channelExists {
					channel.Clients[client.Userid] = client
				}else{
					channel = Channel{
						Clients: map[int]Client{client.Userid: client},
						ChannelId: channelId,
					}
					chatserver.channels[channelId] = channel
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
					chatserver.db.Where("channelid = ?",channel.ChannelId).Order("messagetime DESC").Limit(msgNum).Offset(0).Find(&msgs)
					
					var chatMsgs []model.ChatMessage
					for i := range msgs {
						var user model.User
						chatserver.db.First(&user, msgs[i].Userid)
						chatMsgs = append(chatMsgs, model.ChatMessage{
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
					chatserver.messageQ <- &msg
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

				channel, exists := chatserver.channels[channelId]
				if exists {
					delete(channel.Clients, client.Userid)
				}

				log.Printf("Client %d left channel %d",client.Userid,channelId)

				return
			}
		}
	}
}

func (chatserver *Chatserver) channelLoop() {
	for {
		msg := <- chatserver.messageQ
		// persist
		go chatserver.db.Create(&msg)
		// send to chat channel
		channel := chatserver.channels[msg.Channelid]
		clients := channel.Clients
		for clientId := range clients {
			conn := clients[clientId].Conn
			
			var user model.User
			chatserver.db.First(&user, msg.Userid)
			b, err := json.Marshal(model.ChatMessage{
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

func (chatserver *Chatserver) Run() {
	go chatserver.channelLoop()

	handler := chatserver.createHandler()
	http.HandleFunc("/ws", handler)
	addr := fmt.Sprintf("%s:%s", chatserver.serverHost, chatserver.serverPort)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Listening on address", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewChatserver(serverHost string, serverPort string, db *gorm.DB, jwtKeyFilename string) (Chatserver, error){
	return Chatserver{
		serverHost: serverHost,
		serverPort: serverPort,
		db: db,
		messageQ: make(chan *Message),
		channels: make(map[int]Channel),
		pubKey: utils.ReadJWTKey(jwtKeyFilename),
	}, nil
}