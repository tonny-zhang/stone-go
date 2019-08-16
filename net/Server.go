package net

import (
	"crypto/tls"
	"net"
)

// Server tcp server
type Server struct {
	address        string // Address to open connection: localhost:9999
	config         *tls.Config
	onListen       func()
	onError        func(e error)
	onNewClient    func(c *Client)
	onClientClosed func(c *Client)
	onMessage      func(c *Client, code int16, message string)
}

// OnNewClient Called right after server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClient = callback
}

// OnClientClosed Called right after connection closed
func (s *Server) OnClientClosed(callback func(c *Client)) {
	s.onClientClosed = callback
}

// OnMessage Called when Client receives new message
func (s *Server) OnMessage(callback func(c *Client, code int16, message string)) {
	s.onMessage = callback
}

// OnListen emit on server listen
func (s *Server) OnListen(callback func()) {
	s.onListen = callback
}

// OnError emit on server error
func (s *Server) OnError(callback func(e error)) {
	s.onError = callback
}

// Address get address
func (s *Server) Address() string {
	return s.address
}

// Listen starts network server
func (s *Server) Listen(address string) {
	s.address = address
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		if s.onError != nil {
			s.onError(err)
		}
	}
	defer listener.Close()

	if s.onListen != nil {
		s.onListen()
	}
	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn: conn,
		}
		client.OnClose(func() {
			s.onClientClosed(client)
		})
		client.OnMessage(func(code int16, message string) {
			s.onMessage(client, code, message)
		})
		client.OnConnect(func() {
			s.onNewClient(client)
		})
		go client.work()
	}
}
