package main

import (
	"fmt"

	"sutext.github.io/entry/client"
)

func main() {
	conf := client.NewConfig()
	c := client.New(conf)
	// ctx := context.Background()
	c.Connect("user1", "access_token1")
	sendText(c)
	// for {
	// 	select {
	// 	case status := <-c.NotifyStatus:
	// 		fmt.Println("Received status:", status)
	// 		switch status {
	// 		case client.StatusOpened:
	// 			sendText(c)
	// 		case client.StatusClosed:
	// 			return
	// 		default:
	// 		}
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

func sendText(c *client.Client) {
	var id int64 = 0
	for {
		id++
		var text string
		fmt.Println("Please input text or 'exit' to exit")

		fmt.Scan(&text)
		c.SendText(text, id)
		if text == "exit" {
			return
		}
	}
}
