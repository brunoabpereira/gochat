package main

import (
	"log"
	"os"
	"strconv"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
)

type Payload struct {
	Op string
	Value string
}

func client(v string){
	requestHeader := http.Header{}
	cookie := http.Cookie{Name: "userid", Value: v, Path: "/", HttpOnly: true, Secure: false}
	requestHeader.Set("Cookie", cookie.String())
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", requestHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// join channel
	p := Payload{Op: "join", Value: "1"}
	err = conn.WriteJSON(p)
	if err != nil {
		log.Println(err)
	}
	
	// write every 5 seconds
	go func(){
		i := 0
		for {
			time.Sleep(time.Second*5)

			p := Payload{Op: "send", Value: "Test " + v + " " + strconv.Itoa(i)}
			err = conn.WriteJSON(p)
			if err != nil {
				log.Println(err)
				return
			}

			i++
		}
	}()
	
	// read as fast as possible
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Received message:", string(p))
	}
}

func main(){
	args := os.Args[1:]
	client(args[0])
}