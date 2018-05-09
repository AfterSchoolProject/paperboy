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

type Channel struct {
	Name string `json:"name"`
}

type Subscriber struct {
	Name        string `json:"name"`
	ChannelName string `json:"channelName"`
}

type ChannelMessage struct {
	Channel string `json:"channel"`
	Body    string `json:"body"`
}

var (
	database = redis.DialDatabase(0)
)

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379", database)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Server not available")
		return
	}

	var channel Channel

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &channel)
	if err != nil {
		fmt.Fprintf(w, "unprocessable entity")
		return
	}

	_, err = conn.Do("SET", channel.Name, 1)
	if err != nil {
		fmt.Fprintf(w, "channel not created")
		return
	}

	fmt.Fprintf(w, "channel created!")
}

func Publish(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379", database)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Server not available")
	}

	var message ChannelMessage

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &message)
	if err != nil {
		fmt.Fprintf(w, "unprocessable entity")
	}

	available := conn.Do("GET", message.Channel)
	if available == nil {
		fmt.Fprintf(w, "channel not yet created")
	} else {
		conn.Do("SADD", message.Channel, message.Body)
		fmt.Fprintf(w, "message published")
	}
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", ":6379", database)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Server not available")
		return
	}

	var subscriber Subscriber

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &subscriber)
	if err != nil {
		fmt.Fprintf(w, "unprocessable entity")
	}

	available := conn.Do("GET", subscriber.ChannelName)
	if available == nil {
		fmt.Fprintf(w, "channel not yet created")
		return
	}

	key = Sprintf("%s:subscribers", subscriber.ChannelName)
	conn.Do("SADD", key, subscriber.Name)

	fmt.Fprintf(w, "subscribe")
}

func GetMessage(w http.ResponseWriter, _ *http.Request) {
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/channel", CreateChannel).Methods("POST")
	router.HandleFunc("/publish", Publish).Methods("POST")
	router.HandleFunc("/subscribe", Subscribe).Methods("POST")
	router.HandleFunc("/subscribe/{channel}/message", GetMessage).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
