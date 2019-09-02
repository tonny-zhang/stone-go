package service

import (
	"encoding/json"
	"fmt"
	"stone/logger"
	"stone/net"
	"sync"
	"testing"
)

var conf = Conf{"127.0.0.1", 7000, "../cache", "hello"}
var isStarted = false

func startService() {
	if !isStarted {
		go Start(conf)
		isStarted = true
	}
}
func Test1(t *testing.T) {
	chGet1 := make(chan string)
	chGet2 := make(chan string)
	startService()

	go func() {
		loggerClient := logger.GetPrefixLogger("client1")
		client := net.NewClient()
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

func Test2(t *testing.T) {
	var waitgroup sync.WaitGroup
	startService()

	go func() {
		loggerClient := logger.GetPrefixLogger("client2.1")
		client := &net.Client{}
		client.OnConnError(func(e error) {
			t.Error(e)
		})
		client.OnConnect(func() {
			client.SendMsg(net.CodeBind, conf.Secret)
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": 283, \"c\": \"cb_21\"}")
		})
		client.OnMessage(func(code int16, message string) {
			loggerClient.PrintInfof("code = %d, message = %s", code, message)
			if net.CodeGet == code {
				client.Close()
				waitgroup.Done()
			}
		})
		client.Conn(conf.Host, conf.Port)
	}()

	go func() {
		loggerClient := logger.GetPrefixLogger("client2.2")
		client := &net.Client{}
		client.OnConnError(func(e error) {
			t.Error(e)
		})
		client.OnConnect(func() {
			client.SendMsg(net.CodeBind, conf.Secret)
			client.SendMsg(net.CodeGet, "{\"i\": 1, \"k\": 283, \"c\": \"cb_22\"}")
		})
		client.OnMessage(func(code int16, message string) {
			loggerClient.PrintInfof("code = %d, message = %s", code, message)
			if net.CodeGet == code {
				client.Close()
				waitgroup.Done()
			}
		})
		client.Conn(conf.Host, conf.Port)
	}()
	waitgroup.Add(1)
	waitgroup.Add(1)
	waitgroup.Wait()
	// t.SkipNow()
}
