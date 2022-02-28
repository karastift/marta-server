package world

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const ClientPort = 2222

type Client struct {
	Conn         net.Conn
	LastActive   time.Time
	Port         int
	TimeoutCount int
}

func NewClient(conn net.Conn) *Client {
	client := Client{
		Conn:         conn,
		LastActive:   time.Now(),
		Port:         ClientPort,
		TimeoutCount: 0,
	}

	return &client
}

func (client *Client) Send(data []byte) ([]byte, error) {

	// TODO: fix this method whatever is responsable for the not receiving a response from the client

	// send data
	client.Conn.Write(data)

	// set deadline to 5 seconds
	// so the server doesnt halt
	client.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// read response until newline
	responseData, err := bufio.NewReader(client.Conn).ReadBytes('\n')

	if err != nil {
		client.TimeoutCount++
	}

	return responseData, err
}

func (client *Client) Equals(other Client) bool {
	return client.Conn.RemoteAddr() == other.Conn.RemoteAddr()
}

func (client *Client) UpdateTime() {
	client.LastActive = time.Now()
}

func (client *Client) String() string {
	return fmt.Sprintf("Client(Port: %d, LastActive: %s, Conn: [not printable])", client.Port, client.LastActive.String())
}
