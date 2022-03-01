package world

import (
	"sync"
)

type Clients struct {
	ClientsArray []Client
}

// Returns pointer to an instance of Clients.
func NewClients() *Clients {
	clients := Clients{
		ClientsArray: make([]Client, 0),
	}

	return &clients
}

// Adds `newClient` to the `ClientsArray`.
func (clients *Clients) AddClient(newClient Client) {
	clients.ClientsArray = append(clients.ClientsArray, newClient)
}

// Removes `toRemove` (instance of Client) of `ClientsArray`.
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

// Send data to all clients without waiting for any response.
func (clients *Clients) Send(data string) {

	// waitgroup for waiting until data is sent to all clients
	var wg sync.WaitGroup

	for _, client := range clients.ClientsArray {

		// increment waitgroup counter
		wg.Add(1)

		// start goroutine to send to clients concurrently
		go func(curr Client) {
			// decrement waitgroup counter
			defer wg.Done()

			// send without waiting for response
			err := curr.Send(data)

			if err != nil {
				panic(err)
			}

		}(client)

		// wait until counter is zero
		wg.Wait()
	}
}

// Send data to all clients and wait for response.
//
// Uses the `Send(string)` method of `Client`. This results in increasing the `TimeoutCount` of a client if it doesnt repond.
//
// After the data got sent to a client the method checks if `TimeoutCount` is bigger than 4 and if it is, the client gets removed from `ClientsArray`.
func (clients *Clients) SendWithRes(data string) []string {

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
			res, err := curr.SendWithRes(data)

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
