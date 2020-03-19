package main

import (
	"flag"
	// "go-test-2/errorMyself/alarm"
	"stone/logger"
	"stone/net"
	"stone/service"
)

var (
	host   string
	port   int
	dbPath string
	secret string
	help   bool
)

func testFlag() {
	flag.StringVar(&host, "host", "0.0.0.0", "host for service")
	flag.IntVar(&port, "port", 6006, "port for service")
	flag.StringVar(&dbPath, "dbPath", "./cache", "db path for service")
	flag.StringVar(&secret, "secret", "", "secret for service")
	flag.BoolVar(&help, "help", false, "for help")

	flag.Parse()

	if help {
		flag.PrintDefaults()
	} else {
		testServer()
	}
}
func testServer() {
	var conf = service.Conf{
		Host:   host,
		Port:   port,
		DbPath: dbPath,
		Secret: secret,
	}
	go func() {
		service.Start(conf, func() {
			go func() {
				loggerClient := logger.GetPrefixLogger("client")
				client := net.NewClient()
				client.OnConnError(func(e error) {
					loggerClient.PrintError(e)
				})
				client.OnConnect(func() {
					loggerClient.PrintInfof("%s connect", client.GetConn().LocalAddr())
					client.SendMsg(net.CodeBind, conf.Secret)
					client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": 283, \"c\": \"cb_1\"}")
					client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": -100, \"c\": \"cb_2\"}")
				})
				client.OnMessage(func(code int16, message string) {
					loggerClient.PrintInfof("client get code: %d, message: %s", code, message)
				})
				client.Conn(conf.Host, conf.Port)
			}()
		})
	}()

	select {}
}
func testLoadOtherMod() {
	// alarm.GetError()
}
func testLogger() {
	loggerTest := logger.GetPrefixLogger("test")
	loggerTest.PrintInfof("%s %d", "hello", 123)
}
func main() {
	testLoadOtherMod()
	testFlag()
	// testServer()
	// testLogger()
}
