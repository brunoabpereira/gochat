package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func userIdFromCookie(cookie http.Cookie) int{
	val, _ := strconv.Atoi(cookie.Value)
	return val
}

func createHandler(channels *map[int]Channel, messageQ chan<- *Message, upgrader websocket.Upgrader, db *gorm.DB) func(http.ResponseWriter,*http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
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
		
		// cookie, err := r.Cookie("userid")
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// var client Client = Client{Userid: userIdFromCookie(*cookie)}
		var client Client = Client{Userid: 2}
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

				val, exists := (*channels)[channelId]
				channel = val
				client.Conn = conn
				if exists {
					// _, exists := channel.Clients[client.Userid]
					// if exists {
					// 	log.Printf("Client %d already in Channel %d",client.Userid,channel.ChannelId)
					// 	return
					// }else {
					// 	channel.Clients[client.Userid] = client
					// }

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
		go writeToDb(db, msg)
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

func writeToDb(db *gorm.DB, msg *Message){
	db.Create(&msg)
}

func main() {
	var upgrader = websocket.Upgrader{ReadBufferSize:  1024, WriteBufferSize: 1024}
	var connList map[int]Channel = make(map[int]Channel)
	var messageQ chan *Message = make(chan *Message)
	
	var dsn string = "host=localhost user=gochat password=gochat dbname=gochat port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	go channelLoop(db, &connList, messageQ)

	handler := createHandler(&connList, messageQ, upgrader, db)
	http.HandleFunc("/ws", handler)
	server := &http.Server{
		Addr:              "localhost:9000",
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Println("Listening...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}