package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

type ApiError struct {
	Message string `json:"message"`
}

type ListRes struct {
	Data  []Client `json:"data"`
	Error ApiError `json:"error"`
}

type ShellCmdReq struct {
	ClientId string `json:"clientId"`
	Command  string `json:"command"`
}

type ShellCmdRes struct {
	Data  string   `json:"data"`
	Error ApiError `json:"error"`
}

type InitShellReq struct {
	ClientId string `json:"clientId"`
}

type InitShellRes struct {
	Data  bool     `json:"data"`
	Error ApiError `json:"error"`
}

type KickReq struct {
	ClientId string `json:"clientId"`
}

type KickRes struct {
	Data  bool     `json:"data"`
	Error ApiError `json:"error"`
}

type PingReq struct {
	ClientId string `json:"clientId"`
}

type PingRes struct {
	Data  bool     `json:"data"`
	Error ApiError `json:"error"`
}

var upgrader = websocket.Upgrader{}

type Api struct {
	activeShells map[string]*Shell
}

func NewApi() *Api {
	api := Api{
		activeShells: map[string]*Shell{},
	}

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
	http.HandleFunc("/initShell", initShellHandle)
	http.HandleFunc("/shellCmd", shellCmdHandle)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("API_PORT")), nil))

	return nil
}

func pingHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" {
		return
	} else {

		req := PingReq{}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			json.NewEncoder(w).Encode(PingRes{
				Data: false,
				Error: ApiError{
					Message: err.Error(),
				},
			})
			return
		}

		client, err := clients.GetClientById(req.ClientId)

		if err != nil {
			json.NewEncoder(w).Encode(PingRes{
				Data: false,
				Error: ApiError{
					Message: err.Error(),
				},
			})
			return
		}

		ok := client.Ping()

		json.NewEncoder(w).Encode(PingRes{
			Data:  ok,
			Error: ApiError{},
		})
	}

}

func kickHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != "POST" {
		return
	} else {

		req := KickReq{}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			json.NewEncoder(w).Encode(KickRes{
				Data: false,
				Error: ApiError{
					Message: err.Error(),
				},
			})
			return
		}

		client, err := clients.GetClientById(req.ClientId)

		if err != nil {
			json.NewEncoder(w).Encode(KickRes{
				Data: false,
				Error: ApiError{
					Message: err.Error(),
				},
			})
			return
		}

		clients.RemoveClient(*client)

		json.NewEncoder(w).Encode(KickRes{
			Data:  true,
			Error: ApiError{},
		})
	}
}

func listHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(ListRes{
		Data:  clients.GetAllClients(),
		Error: ApiError{},
	})
}

func initShellHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req := InitShellReq{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		json.NewEncoder(w).Encode(InitShellRes{
			Data: false,
			Error: ApiError{
				Message: err.Error(),
			},
		})
		return
	}

	client, err := clients.GetClientById(req.ClientId)

	if err != nil {
		json.NewEncoder(w).Encode(InitShellRes{
			Data: false,
			Error: ApiError{
				Message: "client id does not exist",
			},
		})
		return
	}

	api.activeShells[client.Id] = NewShell(client)

	json.NewEncoder(w).Encode(InitShellRes{
		Data:  true,
		Error: ApiError{},
	})
}

func shellCmdHandle(w http.ResponseWriter, r *http.Request) {

	// allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req := &ShellCmdReq{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		json.NewEncoder(w).Encode(ShellCmdRes{
			Data: "",
			Error: ApiError{
				Message: err.Error(),
			},
		})
		return
	}

	shell, ok := api.activeShells[req.ClientId]

	if !ok {
		json.NewEncoder(w).Encode(ShellCmdRes{
			Data: "",
			Error: ApiError{
				Message: "no shell is active to a client with that id",
			},
		})
		return
	}

	out, err := shell.Exec(req.Command)

	if err != nil {
		json.NewEncoder(w).Encode(ShellCmdRes{
			Data: "",
			Error: ApiError{
				Message: err.Error(),
			},
		})
		return
	}

	json.NewEncoder(w).Encode(ShellCmdRes{
		Data:  string(out),
		Error: ApiError{},
	})
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
