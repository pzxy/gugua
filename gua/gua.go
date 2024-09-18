package main

import (
	"flag"
	"github.com/firstrow/tcp_server"
	"os"
)

var port = flag.String("port", "50508", "Hostname of the server")

func main() {
	flag.Parse()
	server := tcp_server.New(":" + *port)
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
