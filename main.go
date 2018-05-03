package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateChannel(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "channel created!")
}
func Publish(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "published!")
}

func Subscribe(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "subscribe")
}

func GetMessage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "get-subscription")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/channel", CreateChannel).Methods("POST")
	router.HandleFunc("/publish", Publish).Methods("POST")
	router.HandleFunc("/subscribe", Subscribe).Methods("POST")
	router.HandleFunc("/subscribe/{channel}/message", GetMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
