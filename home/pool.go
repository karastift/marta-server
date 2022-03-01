package home

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/karastift/marta-server/world"
)

// types of requests:
// a login in request
// a response (i hope i dont have to handle them in pool)

type Pool struct {
	Running  bool
	Pausing  bool
	listener net.Listener
	pClients *world.Clients
}

func NewPool(pClients *world.Clients) *Pool {
	pool := Pool{
		Running:  false,
		Pausing:  false,
		pClients: pClients,
	}

	return &pool
}

func (pool *Pool) Start() {
	pool.Running = true

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", world.ClientPort))
	pool.listener = listener

	// check if initiating the listener failed
	if err != nil {
		fmt.Println("Failed initiating the listener.")
		fmt.Println(err)
		os.Exit(1)
	}

	// main loop to check of incoming connections
	for {
		if pool.Pausing {
			continue
		}
		// accepting an incoming connection
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		// read stream until \n
		netData, err := bufio.NewReader(conn).ReadString('\n')

		// if reading the incoming data failes, just ignore the connection
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			// a client wants to login
			// append client to clients
			if netData == "marta login\n" {
				client := world.NewClient(conn)
				pool.pClients.AddClient(*client)

				fmt.Println("Found a new client.")

				client.Send([]byte("marta logged in\n"))

			} else {
				fmt.Println("Received in pool: '" + netData + "'")
			}
		}
	}
}

func (pool *Pool) Stop() {
	err := pool.listener.Close()

	if err != nil {
		fmt.Println("Failed to close listener.")
		fmt.Println(err)
	}

	pool.Running = false
}

func (pool *Pool) Pause() {
	pool.Pausing = true
}

func (pool *Pool) Resume() {
	pool.Pausing = false
}
