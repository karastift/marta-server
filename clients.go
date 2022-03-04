package main

import (
	"errors"
	"sync"
)

const maxTimeouts int = 2

type Clients struct {
	clientsMap map[string]Client
}

// Returns pointer to an instance of Clients.
func NewClients() *Clients {
	clients := Clients{
		clientsMap: make(map[string]Client),
	}

	return &clients
}

// Returns client with curresponding id.
func (clients *Clients) GetClientById(id string) (*Client, error) {

	client, ok := clients.clientsMap[id]

	if !ok {
		return nil, errors.New("could not find client: no id matches the given one")

	} else {
		return &client, nil
	}
}

// Returns array of clients.
func (clients *Clients) GetAllClients() []Client {
	arr := make([]Client, 0, len(clients.clientsMap))

	for _, client := range clients.clientsMap {
		arr = append(arr, client)
	}

	return arr
}

// Adds `newClient` to the `ClientsArray`.
func (clients *Clients) AddClient(newClient Client) {
	clients.clientsMap[newClient.Id] = newClient
}

// Removes `toRemove` (instance of Client) of `ClientsArray`.
func (clients *Clients) RemoveClient(toRemove Client) {
	delete(clients.clientsMap, toRemove.Id)
}

// Pings all clients and returns map: client -> ok
func (clients *Clients) Ping() map[Client]bool {

	oks := map[Client]bool{}

	mapWithWaitGroup(&clients.clientsMap, func(curr *Client, wg *sync.WaitGroup) {
		// decrement waitgroup counter
		defer wg.Done()

		oks[*curr] = curr.Ping()
	})

	return oks
}

// Send data to all clients without waiting for any response.
func (clients *Clients) Send(data string) {

	mapWithWaitGroup(&clients.clientsMap, func(curr *Client, wg *sync.WaitGroup) {
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
// After the data got sent to a client the method checks if `TimeoutCount` is bigger than `maxTimeouts` and if it is, the client gets removed from `ClientsArray`.
func (clients *Clients) SendWithRes(data string) []string {

	// array that gatheres all responses
	responses := make([]string, len(clients.clientsMap))
	timedOut := make([]Client, 0)

	mapWithWaitGroup(&clients.clientsMap, func(curr *Client, wg *sync.WaitGroup) {

		// decrement waitgroup counter
		defer wg.Done()

		// send and get response
		res, err := curr.SendWithRes(data)

		if err != nil {
			// if count is to big append to timedOut slice
			if curr.TimeoutCount > maxTimeouts {
				timedOut = append(timedOut, *curr)
			}
		} else {
			responses = append(responses, string(res))
		}

		clients.clientsMap[curr.Id] = *curr
	})

	// remove timed out clients from clientarray
	for _, client := range timedOut {
		clients.RemoveClient(client)
	}

	return responses
}

// Requests info from client. Updates `Info` and returns it.
func (clients *Clients) RequestInfo() {

	mapWithWaitGroup(&clients.clientsMap, func(curr *Client, wg *sync.WaitGroup) {
		defer wg.Done()
		curr.RequestInfo()

		// client changed, so change it in array too
		clients.clientsMap[curr.Id] = *curr
	})
}

// Starts the function as a goroutine for every client in the array.
// The function has to defer call `wg.Done()` when its finished!
func mapWithWaitGroup(arr *map[string]Client, fn func(*Client, *sync.WaitGroup)) {

	// create waitgroup
	var wg sync.WaitGroup

	for _, client := range *arr {

		// increment waitgroup counter
		wg.Add(1)

		// call passed function as goroutine
		go fn(&client, &wg)

		// wait until counter is zero
		wg.Wait()
	}
}
