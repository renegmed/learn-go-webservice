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
	go func(c *websocket.Conn) {
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("received message %s\n", msg.Data)
		}
	}(ws)

	for {
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
		/*
			websocket has not properly closed when the browser is terminated
			thus we are getting console error when the next reading comes
				2020/06/19 16:59:29 EOF
				2020/06/19 16:59:49 write tcp [::1]:5000->[::1]:44712: write: broken pipe
		*/
	}
}
