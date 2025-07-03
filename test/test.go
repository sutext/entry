package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"sutext.github.io/entry/backoff"
	"sutext.github.io/entry/client"
	"sutext.github.io/entry/packet"
)

type Client struct {
	cli     *client.Client
	userId  string
	token   string
	backoff backoff.Backoff
	count   int
}

func RandomClient() *Client {
	return &Client{
		cli:     client.New(client.NewConfig()),
		userId:  fmt.Sprintf("user_%d", rand.Int()),
		token:   fmt.Sprintf("access_token_%d", rand.Int()),
		backoff: backoff.Random(time.Second*5, time.Second*10),
		count:   0,
	}
}
func (c *Client) Start() {
	c.cli.Connect(&packet.Identity{Token: c.userId, UserID: c.userId, ClientID: c.token})
	for {
		c.cli.SendText("hello world")
		time.Sleep(c.backoff.Next(c.count))
		c.count++
	}
}
func addClient(count uint) {
	for range count {
		c := RandomClient()
		go c.Start()
	}
}
func main() {
	go func() {
		timer := time.NewTimer(5 * time.Second)
		ctx := context.Background()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				addClient(10)
				timer.Reset(5 * time.Second)
			}
		}
	}()
	select {}
}
