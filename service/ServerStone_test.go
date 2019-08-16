package service

import (
	"encoding/json"
	"fmt"
	"stone/logger"
	"stone/net"
	"testing"
)

var conf = Conf{"127.0.0.1", 7000, "../cache", "hello"}

func startService() {
	Start(conf)
}
func Test1(t *testing.T) {
	chGet1 := make(chan string)
	chGet2 := make(chan string)
	go startService()

	go func() {
		loggerClient := logger.GetPrefixLogger("client")
		client := &net.Client{}
		client.OnConnError(func(e error) {
			t.Error(e)
		})
		client.OnConnect(func() {
			client.SendMsg(net.CodeBind, conf.Secret)
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": 283, \"c\": \"cb_1\"}")
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": -100, \"c\": \"cb_2\"}")
		})
		client.OnMessage(func(code int16, message string) {
			loggerClient.PrintInfof("code = %d, message = %s", code, message)
			if net.CodeGet == code {
				data := make(map[string]interface{})
				json.Unmarshal([]byte(message), &data)
				// fmt.Println(data)
				cbName := data["c"].(string)
				// fmt.Println("cbName = ", cbName)
				if cbName == "cb_1" {
					chGet1 <- cbName
				} else if cbName == "cb_2" {
					chGet2 <- cbName
				}
			} else if net.CodeError == code {
				data := make(map[string]interface{})
				json.Unmarshal([]byte(message), &data)
				fmt.Println(data)
			}
		})
		client.Conn(conf.Host, conf.Port)
	}()

	isGet1 := false
	isGet2 := false
	fnTest := func() {
		if isGet1 && isGet2 {
			// fmt.Println("all down")
			t.SkipNow()
		}
	}
	for {
		select {
		case <-chGet1:
			isGet1 = true
			fnTest()
		case <-chGet2:
			isGet2 = true
			fnTest()
		}
	}

}
