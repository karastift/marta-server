package world

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn         net.Conn
	TimeoutCount int
}

// Returns a pointer to an instance of Client.
func NewClient(conn net.Conn) *Client {
	client := Client{
		Conn: conn,
	}

	return &client
}

// Send data to client without waiting for any response.
func (client *Client) Send(data string) error {

	_, err := client.Conn.Write([]byte(data))

	return err
}

// Send data to client and wait for response.
//
// Increases `TimeoutCount` of client if it doesnt respond after 5 seconds.
func (client *Client) SendWithRes(data string) (string, error) {

	client.Conn.Write([]byte(data))

	// set deadline to in 5 seconds
	// if client does not respond after 5 seconds, it resumes and the timeout counter get incremented
	client.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	res, err := bufio.NewReader(client.Conn).ReadString('\n')

	if err != nil {
		client.TimeoutCount++
	}

	return res, err
}

// Checks if client is equal to the other based on there local address.
func (client *Client) Equals(other Client) bool {
	return client.Conn.LocalAddr() == other.Conn.LocalAddr()
}

// Returns string representation of itsself.
func (client *Client) String() string {
	return fmt.Sprintf("Client(%s)", client.Conn.LocalAddr())
}
