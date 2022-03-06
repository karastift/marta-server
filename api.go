package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{}

type Api struct{}

func NewApi() *Api {
	api := Api{}

	return &api
}

func (api *Api) Serve() error {

	// let functions return output and then write it to connection
	// http endpoint for basic commands
	// for shell, i can try using websockets

	http.HandleFunc("/ws", wsHandle)
	http.HandleFunc("/list", listHandle)
	http.HandleFunc("/kick", kickHandle)
	http.HandleFunc("/ping", pingHandle)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("API_PORT")), nil))

	return nil
}

func pingHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	} else {

		// allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		clientId, err := bufio.NewReader(r.Body).ReadString('\n')

		if err != nil {
			json.NewEncoder(w).Encode(false)
			return
		}
		clientId = strings.TrimSpace(clientId)

		client, err := clients.GetClientById(clientId)

		if err != nil {
			json.NewEncoder(w).Encode(false)
			return
		}

		ok := client.Ping()

		json.NewEncoder(w).Encode(ok)
	}

}

func kickHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	} else {

		// allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		clientId, err := bufio.NewReader(r.Body).ReadString('\n')

		if err != nil {
			json.NewEncoder(w).Encode(false)
			return
		}
		clientId = strings.TrimSpace(clientId)

		client, err := clients.GetClientById(clientId)

		if err != nil {
			json.NewEncoder(w).Encode(false)
			return
		}

		clients.RemoveClient(*client)
		json.NewEncoder(w).Encode(true)
	}
}

func listHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(clients.GetAllClients())
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()

	fmt.Println("Connected!")

	for {
		var message Message

		err := ws.ReadJSON(&message)

		if err != nil {
			fmt.Printf("error occurred: %v", err)
			break
		}

		fmt.Println(message)

		// send message from server
		if err := ws.WriteJSON(message); !errors.Is(err, nil) {
			fmt.Printf("error occurred: %v", err)
		}
	}
}
