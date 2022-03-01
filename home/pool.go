package home

import (
	"bufio"
	"fmt"
	"net"

	"github.com/karastift/marta-server/world"
)

type Pool struct {
	Port     int
	Running  bool
	Pausing  bool
	listener net.Listener
	pClients *world.Clients
}

// Returns pointer to an instance of Pool.
func NewPool(pClients *world.Clients, port int) *Pool {
	pool := Pool{
		Running:  false,
		pClients: pClients,
		Port:     port,
	}

	return &pool
}

// Starts the pool. Pool listens now to incoming tcp connections and adds new clients to pClients.
func (pool *Pool) Start() error {
	pool.Running = true

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", pool.Port))
	pool.listener = listener

	// check if initiating the listener failed
	if err != nil {
		return err
	}

	// main loop to check of incoming connections
	for {
		// accepting an incoming connection
		conn, err := listener.Accept()

		// if accepting the connection failes, just ignore it
		if err != nil {
			continue
		}

		// read stream until \n
		netData, err := bufio.NewReader(conn).ReadString('\n')

		// if reading the incoming data failes, just ignore the connection
		if err != nil {
			continue

		} else {
			// a client wants to login
			// append client to clients
			if netData == "marta login\n" {

				client := world.NewClient(conn)

				client.Send("marta logged in\n")

				pool.pClients.AddClient(*client)

				Log("Client connected: " + client.String())

			} else {
				fmt.Println("Received in pool: '" + netData + "'")
			}
		}
	}
}

// Stop the pool.
func (pool *Pool) Stop() error {
	err := pool.listener.Close()
	pool.Running = false

	if err != nil {
		return err
	}
	return nil
}
