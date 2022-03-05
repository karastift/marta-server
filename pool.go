package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Pool struct {
	Port     int
	Running  bool
	Pausing  bool
	listener net.Listener
}

// Returns pointer to an instance of Pool.
func NewPool(port int) *Pool {
	pool := Pool{
		Running: false,
		Port:    port,
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
			if strings.HasPrefix(netData, "marta login|") {

				client := NewClient(conn)

				clientInfo, err := NewClientInfo(strings.TrimPrefix(netData, "marta login|"))

				if err != nil {
					logger.Println(fmt.Sprintf("Invalid login message: '%s'", netData))
					continue
				}

				client.Info = *clientInfo

				client.Send("marta logged in\n")

				clients.AddClient(*client)

				logger.Println("Client connected: " + client.String())

			} else {
				logger.Println("Received in pool (should not happen): '" + netData + "'")
			}
		}
	}
}

// Stop the pool.
func (pool *Pool) Stop() error {
	err := pool.listener.Close()
	pool.Running = false

	return err
}
