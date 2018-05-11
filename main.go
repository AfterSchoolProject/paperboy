package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// type Channel struct {
// 	Name string `json:"name"`
// }

// type Subscriber struct {
// 	Name        string `json:"name"`
// 	ChannelName string `json:"channelName"`
// }

type ChannelMessage struct {
	Channel string `json:"channel"`
	Body    string `json:"body"`
}

var (
	database = redis.DialDatabase(0)
)

// TODO Add channel to db
// func CreateChannel(w http.ResponseWriter, r *http.Request) {
// 	conn, err := redis.Dial("tcp", ":6379", database)
// 	defer conn.Close()
// 	if err != nil {
// 		fmt.Fprintf(w, "Server not available")
// 		return
// 	}

// 	var channel Channel

// 	body, err := ioutil.ReadAll(r.Body)
// 	err = json.Unmarshal(body, &channel)
// 	if err != nil {
// 		fmt.Fprintf(w, "unprocessable entity")
// 		return
// 	}

// 	_, err = conn.Do("SET", channel.Name, 1)
// 	if err != nil {
// 		fmt.Fprintf(w, "channel not created")
// 		return
// 	}

// 	fmt.Fprintf(w, "channel created!")
// }

func Publish(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379", database)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Server not available")
		return
	}

	var cMsg ChannelMessage

	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &cMsg)
	if err != nil {
		fmt.Fprintf(w, "unprocessable entity")
		return
	}

	conn.Do("SADD", cMsg.Channel, cMsg.Body)
	fmt.Fprintf(w, "message published")
}

// TODO add subscriber to db and add relationship to channel
// func Subscribe(w http.ResponseWriter, r *http.Request) {
// 	conn, err := redis.Dial("tcp", ":6379", database)
// 	defer conn.Close()
// 	if err != nil {
// 		fmt.Fprintf(w, "Server not available")
// 		return
// 	}

// 	var subscriber Subscriber

// 	body, _ := ioutil.ReadAll(r.Body)
// 	err = json.Unmarshal(body, &subscriber)
// 	if err != nil {
// 		fmt.Fprintf(w, "unprocessable entity")
// 		return
// 	}

// 	available, _ := conn.Do("GET", subscriber.ChannelName)
// 	if available == nil {
// 		fmt.Fprintf(w, "channel not yet created")
// 		return
// 	}

// 	key := fmt.Sprintf("%s:subscribers", subscriber.ChannelName)
// 	conn.Do("SADD", key, subscriber.Name)

// 	fmt.Fprintf(w, "subscribed")
// }

func GetMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379", database)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Server not available")
	}

	vars := mux.Vars(r)
	message, _ := redis.String(conn.Do("SPOP", vars["channel"]))

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/channel", CreateChannel).Methods("POST")
	router.HandleFunc("/publish", Publish).Methods("POST")
	router.HandleFunc("/subscribe", Subscribe).Methods("POST")
	router.HandleFunc("/subscribe/{channel}/message", GetMessage).Methods("GET")

	fmt.Println("Starting server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
