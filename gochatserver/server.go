package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
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

func createHandler(channels *map[int]Channel, messageQ chan<- *Message, upgrader websocket.Upgrader) func(http.ResponseWriter,*http.Request){
	return func (w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		log.Println("New connection from",conn.RemoteAddr().String())
		
		cookie, err := r.Cookie("userid")
		if err != nil {
			log.Println(err)
			return
		}
		var client Client = Client{Userid: userIdFromCookie(*cookie)}
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
					_, exists := channel.Clients[client.Userid]
					if exists {
						log.Printf("Client %d already in Channel %d",client.Userid,channel.ChannelId)
						return
					}else {
						channel.Clients[client.Userid] = client
					}
				}else{
					channel = Channel{
						Clients: map[int]Client{client.Userid: client},
						ChannelId: channelId,
					}
					(*channels)[channelId] = channel
				}

				inChannel = true
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
			err := (*conn).WriteMessage(websocket.TextMessage, []byte(msg.Messagetext))
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

	handler := createHandler(&connList, messageQ, upgrader)
	http.HandleFunc("/ws", handler)
	server := &http.Server{
		Addr:              "localhost:8080",
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Println("Listening...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}