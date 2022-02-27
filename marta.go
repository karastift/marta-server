package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {

	port := ":" + "2222"

	l, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	conn, err := l.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"

		fmt.Print("-> ", string(netData))
		conn.Write([]byte(myTime))
	}
}
