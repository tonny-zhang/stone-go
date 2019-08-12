package main

import (
	"log"
	"stone/net"
)

func main() {
	server := &net.Server{}

	server.OnListen(func(s *net.Server) {
		log.Printf("listen at %s\n", s.Address())
	})
	server.OnNewClient(func(c *net.Client) {
		log.Printf("client [%s] connect \n", c.Conn().RemoteAddr())
	})
	server.OnNewMessage(func(c *net.Client, code int16, message string) {
		log.Printf("get client [%s] code [%d] message [%s] \n", c.Conn().RemoteAddr(), code, message)
		c.Send([]byte("from server get ["+message+"]"), -100)
	})
	server.OnClientClosed(func(c *net.Client, e error) {
		log.Printf("client [%s] closed\n", c.Conn().RemoteAddr())
	})
	server.Listen("0.0.0.0:6006")
}
