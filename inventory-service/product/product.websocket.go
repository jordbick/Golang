package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// Message type
type message struct {
	Data string `json: "data"`
	Type string `json: "type"`
}

// websocket handler that accepts a websocket connection
func productSocket(ws *websocket.Conn) {
	// Need to use channel to communicate with our Go routine. Need to close our connections once our client disconnects
	// Use channel to signal to handler the connection is closed
	done := make(chan struct{})
	fmt.Println("new websocket connection established")

	// Listen for incoming data on the WebSocket in a go routine
	// for loop, which will run forever
	go func(c *websocket.Conn) {
		for {
			var msg message
			// codec JSON type to call Receive. Pass in message struct to retrieve data from the client
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			// If we get any data it'll print out what was received
			fmt.Printf("received message %s\n", msg.Data)
		}
		close(done)

	}(ws)

	// Simulate the product data being updated by running a query every 10 seconds
	// Query to get top 10 products and write the list to the WebSocket
loop:
	for {
		// select case statement to listen for closing of channel
		select {
		case <-done:
			fmt.Println("connection was closed")
			break loop
		default:
			// This function in product.data file
			products, err := GetTopTenProducts()
			if err != nil {
				log.Println(err)
				break
			}
			// websocket.JSON.Send method to handle the marshalling of our slice of products into JSON
			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Println(err)
				break
			}
			time.Sleep(10 * time.Second)
		}

	}
	fmt.Println("Closing websocket")
	defer ws.Close()
}

// Also need to setup handler,  add this to the setupRoutes function in the product.service file
