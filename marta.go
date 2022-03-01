// server

package main

import (
	"github.com/karastift/marta-server/home"
	"github.com/karastift/marta-server/world"
)

var clients = world.NewClients()

func main() {

	pool := home.NewPool(clients, 2222)

	pool.Start()

}
