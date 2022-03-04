package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const timeOutDuration int = 5

type Client struct {
	Id           string
	Conn         net.Conn
	TimeoutCount int
	Info         ClientInfo
}

// Returns a pointer to an instance of Client.
func NewClient(conn net.Conn) *Client {
	client := Client{
		Id:   randStringBytes(5),
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
	client.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeOutDuration) * time.Second))

	res, err := bufio.NewReader(client.Conn).ReadString('\n')

	if err != nil {
		client.TimeoutCount = client.TimeoutCount + 1
	}

	return res, err
}

// Requests info from client. Updates `Info` and returns it.
func (client *Client) RequestInfo() (*ClientInfo, error) {
	res, err := client.SendWithRes("!info\n")

	if err != nil {
		return nil, err
	}

	info, err := NewClientInfo(res)

	if err != nil {
		return nil, err
	}

	client.Info = *info

	return info, nil
}

// Ping the client. Returns true if client responded.
func (client *Client) Ping() bool {
	res, err := client.SendWithRes("!ping\n")

	return err == nil && res == "pong\n"
}

// Checks if client is equal to the other based on there local address.
func (client *Client) Equals(other Client) bool {
	return client.Id == other.Id && client.Conn.LocalAddr() == other.Conn.LocalAddr()
}

// Returns string representation of itsself.
func (client *Client) String() string {
	return fmt.Sprintf("Client(Id: %s, Addr: %s)", client.Id, client.Conn.LocalAddr())
}

// Returns a random string with the given length.
func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
