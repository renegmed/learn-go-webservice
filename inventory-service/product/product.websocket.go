package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	// we can verify that the origin is an allowed origin
	fmt.Printf("origin: %s\n", ws.Config().Origin)

	defer ws.Close()

	done := make(chan struct{})

	go func(c *websocket.Conn) {
		defer close(done)
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("received message %s\n", msg.Data)
		}
	}(ws)
loop:
	for {
		select {
		case <-done:
			log.Println("connection was closed, lets break out of here.")
			break loop
		default:
			log.Println("sending top 10 products list to the client.")
			products, err := GetTopTenProducts()
			if err != nil {
				log.Fatal(err)
			}

			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Println(err)
				break
			}

			// pause for 10 seconds before sending again
			time.Sleep(10 * time.Second)
		}

	}
	fmt.Println("closing the connection")
}
