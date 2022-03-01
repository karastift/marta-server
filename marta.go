// server

package main

import (
	"fmt"

	"github.com/karastift/marta-server/home"
	"github.com/karastift/marta-server/world"
)

var clients = world.NewClients()

func main() {

	// change server port
	// just try sending from client and see if server receives

	pool := home.NewPool(clients)

	go pool.Start()

	fmt.Scanf("%s")

	pool.Pausing = true
	responses := clients.Send("Love you so much, it makes me sick.\n")
	pool.Pausing = false

	for res := range responses {
		fmt.Println(res)
	}

	pool.Stop()
}
