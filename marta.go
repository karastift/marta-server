// server

package main

import (
	"bufio"
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
	// get right macaddress

	// shell access (could be hard)
	// file download (could be hard)
	// maybe update command (could be hard)

	logger = log.Default()

	pool := NewPool(clients, 2222)

	go pool.Start()

	testResponses()
}

func info() {
	clients.RequestInfo()

	for _, client := range clients.ClientsArray {
		fmt.Println(client.Info.String())
	}
}

func testResponses() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		in, _ := reader.ReadString('\n')

		if in == "!info\n" {
			info()
			continue
		}

		responses := clients.SendWithRes(in)

		for _, res := range responses {
			fmt.Println("'" + strings.TrimRight(res, "\n") + "'")
		}

		for _, client := range clients.ClientsArray {
			fmt.Println(client.Conn.RemoteAddr())
		}
	}
}
