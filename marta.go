// server

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var clients = NewClients()
var logger *log.Logger

func main() {

	// TODO:
	// status command on server
	// check out memory modules on ghw
	// check out ghw in general
	// use map for storing clients later (id -> client)
	// handle all errors on server
	// a client can crash, the server shouldnt

	// search for TODO:

	// file download (could be hard)
	// maybe update command (could be hard)

	logger = log.Default()

	pool := NewPool(clients, 2222)

	go pool.Start()

	commandLoop()
}

func info(clientId string) error {

	if clientId != "" {
		client, err := clients.GetClientById(clientId)

		if err != nil {
			return err
		}

		client.RequestInfo()

		fmt.Println(client.Info.String())
	} else {
		clients.RequestInfo()

		for _, client := range clients.ClientsArray {
			fmt.Println(client.Info.String())
		}
	}

	return nil
}

func list() {
	for _, client := range clients.ClientsArray {
		fmt.Println(client.String())
	}
}

func shell(clientId string) error {

	if clientId == "" {
		return errors.New("could not initiate shell: client id was not passed")
	}

	client, err := clients.GetClientById(clientId)

	if err != nil {
		return errors.New("could not initiate shell: client with that id is not in array")
	}

	shell := NewShell(client)

	shell.Start()

	return nil
}

func commandLoop() {
	reader := bufio.NewReader(os.Stdin)
	ctrlDCount := 0

	for {
		fmt.Print("☲ ")

		in, err := reader.ReadString('\n')

		if err != nil {
			ctrlDCount++
			if ctrlDCount > 1 {
				break
			}
			fmt.Println("Press Ctrl + D again, if you want to exit.")
			continue
		}

		in = strings.TrimSpace(in)
		split := strings.Split(in, " ")
		cmd := split[0]
		clientId := ""

		if len(split) > 1 {
			clientId = split[1]
		}

		switch cmd {
		case "!info":
			err := info(clientId)
			if err != nil {
				fmt.Println(err)
			}
		case "!list":
			list()
		case "!shell":
			err := shell(clientId)
			if err != nil {
				fmt.Println(err)
			}
		default:
			responses := clients.SendWithRes(in)

			for _, res := range responses {
				fmt.Println("'" + strings.TrimRight(res, "\n") + "'")
			}
		}
	}
}

// func testResponses() {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("> ")

// 		in, _ := reader.ReadString('\n')

// 		if in == "!info\n" {
// 			info()
// 			continue
// 		}

// 		responses := clients.SendWithRes(in)

// 		for _, res := range responses {
// 			fmt.Println("'" + strings.TrimRight(res, "\n") + "'")
// 		}

// 		for _, client := range clients.ClientsArray {
// 			fmt.Println(client.Conn.RemoteAddr())
// 		}
// 	}
// }
