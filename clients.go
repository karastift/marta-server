package main

import (
	"errors"
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

// Returns client with curresponding id.
func (clients *Clients) GetClientById(id string) (*Client, error) {
	for _, client := range clients.ClientsArray {
		if client.Id == id {
			return &client, nil
		}
	}

	return nil, errors.New("could not find client: no id matches the given one")
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
	mapWithWaitGroup(&clients.ClientsArray, func(_ int, curr *Client, wg *sync.WaitGroup) {
		// decrement waitgroup counter
		defer wg.Done()

		// send without waiting for response
		err := curr.Send(data)

		if err != nil {
			panic(err)
		}
	})
}

// Send data to all clients and wait for response.
//
// Uses the `Send(string)` method of `Client`. This results in increasing the `TimeoutCount` of a client if it doesnt repond.
//
// After the data got sent to a client the method checks if `TimeoutCount` is bigger than 2 and if it is, the client gets removed from `ClientsArray`.
func (clients *Clients) SendWithRes(data string) []string {

	// array that gatheres all responses
	responses := make([]string, len(clients.ClientsArray))
	timedOut := make([]Client, 0)

	mapWithWaitGroup(&clients.ClientsArray, func(i int, curr *Client, wg *sync.WaitGroup) {

		// decrement waitgroup counter
		defer wg.Done()

		// send and get response
		res, err := curr.SendWithRes(data)

		if err != nil {
			// if count is to big append to timedOut slice
			if curr.TimeoutCount > 2 {
				timedOut = append(timedOut, *curr)
			}
		} else {
			responses = append(responses, string(res))
		}

		clients.ClientsArray[i] = *curr
	})

	// remove timed out clients from clientarray
	for _, client := range timedOut {
		clients.RemoveClient(client)
	}

	return responses
}

// Requests info from client. Updates `Info` and returns it.
func (clients *Clients) RequestInfo() {

	mapWithWaitGroup(&clients.ClientsArray, func(index int, curr *Client, wg *sync.WaitGroup) {
		defer wg.Done()
		curr.RequestInfo()

		// client changed, so change it in array too
		clients.ClientsArray[index] = *curr
	})
}

// Starts the function as a goroutine for every client in the array.
// The function has to defer call `wg.Done()` when its finished!
func mapWithWaitGroup(arr *[]Client, fn func(int, *Client, *sync.WaitGroup)) {

	// create waitgroup
	var wg sync.WaitGroup

	for i, client := range *arr {

		// increment waitgroup counter
		wg.Add(1)

		// call passed function as goroutine
		go fn(i, &client, &wg)

		// wait until counter is zero
		wg.Wait()
	}
}