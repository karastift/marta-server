package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Client struct {
	Id   string
	Conn net.Conn
	Info ClientInfo
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

	dur, err := strconv.Atoi(os.Getenv("TIMEOUT_DURATION"))

	if err != nil {
		logger.Panicln(err)
	}

	// set deadline to in `dur` seconds
	client.Conn.SetReadDeadline(time.Now().Add(time.Duration(dur) * time.Second))

	return bufio.NewReader(client.Conn).ReadString('\n')
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

// Returns json representation of itsself as string.
func (client *Client) Json() string {
	j, err := json.Marshal(client)

	if err != nil {
		logger.Panicln(err)
	}

	return string(j)
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
