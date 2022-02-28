// server

package main

import (
	"fmt"

	"github.com/karastift/marta-server/home"
	"github.com/karastift/marta-server/world"
)

var clients = world.NewClients()

func main() {

	pool := home.NewPool(clients)

	go pool.Start()

	fmt.Scanf("%s")

	responses := clients.Send("Love you so much, it makes me sick.\n")

	for res := range responses {
		fmt.Println(res)
	}
}
