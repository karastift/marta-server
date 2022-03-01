// server

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/karastift/marta-server/home"
	"github.com/karastift/marta-server/world"
)

var clients = world.NewClients()

func main() {

	pool := home.NewPool(clients, 2222)

	go pool.Start()

	testResponses()
}

func testResponses() {
	reader := bufio.NewReader(os.Stdin)

	for {
		in, _ := reader.ReadString('\n')

		res := clients.SendWithRes(in)

		fmt.Println(res)
	}
}
