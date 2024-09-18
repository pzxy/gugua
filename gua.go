package main

import (
	"github.com/firstrow/tcp_server"
	"os"
)

func main() {
	server := tcp_server.New("localhost:" + port)
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		if message == "gu\n" {
			name, err := os.Hostname()
			if err != nil {
				_ = c.Send(err.Error())
			} else {
				_ = c.Send("gua:" + name)
			}
		}
	})
	server.Listen()
}
