package main

import (
	"stone/logger"
	"stone/net"
	"stone/service"
)

var conf = service.Conf{
	Host: "0.0.0.0",
	Port: 6006,
}

func testServer() {
	go func() {
		service.Start(conf)
	}()

	go func() {
		loggerClient := logger.GetPrefixLogger("client")
		client := &net.Client{}
		client.OnConnError(func(e error) {
			loggerClient.PrintError(e)
		})
		client.OnConnect(func() {
			loggerClient.PrintInfof("%s connect", client.GetConn().LocalAddr())
			client.SendMsg(net.CodeBind, "hello")
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": 283, \"c\": \"cb_1\"}")
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": -100, \"c\": \"cb_2\"}")
		})
		client.OnMessage(func(code int16, message string) {
			loggerClient.PrintInfof("client get code: %d, message: %s", code, message)
		})
		client.Conn(conf.Host, conf.Port)
	}()

	for {

	}
}
func testLogger() {
	loggerTest := logger.GetPrefixLogger("test")
	loggerTest.PrintInfof("%s %d", "hello", 123)
}
func main() {
	testServer()
	// testLogger()
}
