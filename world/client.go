package world

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn         net.Conn
	LastActive   time.Time
	TimeoutCount int
}

func NewClient(conn net.Conn) *Client {
	client := Client{
		Conn:       conn,
		LastActive: time.Now(),
	}

	return &client
}

func (client *Client) Send(data string) error {

	_, err := client.Conn.Write([]byte(data))

	return err
}

func (client *Client) SendWithRes(data string) (string, error) {

	client.Conn.Write([]byte(data))

	res, err := bufio.NewReader(client.Conn).ReadString('\n')

	return res, err
}

func (client *Client) Equals(other Client) bool {
	return client.Conn.LocalAddr() == other.Conn.LocalAddr()
}

func (client *Client) UpdateTime() {
	client.LastActive = time.Now()
}

func (client *Client) String() string {
	return fmt.Sprintf("Client(%s)", client.Conn.LocalAddr())
}
