package world

import (
	"fmt"
	"sync"
)

type Clients struct {
	ClientsArray []Client
}

func NewClients() *Clients {
	clients := Clients{
		ClientsArray: make([]Client, 0),
	}

	return &clients
}

func (clients *Clients) AddClient(newClient Client) {
	clients.ClientsArray = append(clients.ClientsArray, newClient)
}

func (clients *Clients) RemoveClient(toRemove Client) {

	var removeIndex int

	for index, client := range clients.ClientsArray {
		if client.Equals(toRemove) {
			removeIndex = index
			break
		}
	}

	clients.ClientsArray = append(clients.ClientsArray[:removeIndex], clients.ClientsArray[removeIndex+1:]...)
}

func (clients *Clients) Send(data string) []string {

	// array that gatheres all responses
	responses := make([]string, len(clients.ClientsArray))

	// waitgroup for waiting until the last client reponded or timed out
	var wg sync.WaitGroup

	for _, client := range clients.ClientsArray {

		// increment waitgroup counter
		wg.Add(1)

		// start goroutine to send to clients concurrently
		go func(curr Client) {
			// decrement waitgroup counter
			defer wg.Done()
			// send and get response

			fmt.Println("Sending to: " + curr.String())

			res, err := curr.Send([]byte(data))

			if err != nil {
				// if client has timeouted 4 times now, remove it from connected clients
				if curr.TimeoutCount > 3 {
					// idk if the defer statement is necessary
					defer clients.RemoveClient(curr)
				}
			} else {
				responses = append(responses, string(res))
			}
		}(client)

		// wait until counter is zero
		wg.Wait()
	}

	return responses
}
