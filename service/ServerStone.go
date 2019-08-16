package service

import (
	"encoding/json"
	"fmt"
	"path"
	"stone/logger"
	"stone/net"
	"stone/store"
	"time"
)

// Conf conf for service
type Conf struct {
	Host   string
	Port   int
	DbPath string
	Secret string
}

// Start start service
func Start(conf Conf) {
	server := &net.Server{}
	loggerServer := logger.GetLoggerPrefix("serverStone")
	server.OnListen(func() {
		loggerServer.PrintInfof("listen at %s\n", server.Address())
	})
	server.OnNewClient(func(c *net.Client) {
		loggerServer.PrintInfof("client [%s] connect \n", c.GetKey())

		time.AfterFunc(time.Second*3, func() {
			if !c.GetChecked() {
				loggerServer.PrintInfof("client [%s] close case no bind \n", c.GetKey())
				c.Close()
			}
		})
	})
	server.OnMessage(func(c *net.Client, code int16, message string) {
		loggerServer.PrintInfof("code = %d, message = %s", code, message)
		if code == net.CodeBind {
			if message == conf.Secret {
				c.SetChecked()
				loggerServer.PrintInfof("client [%s] binded\n", c.GetKey())
			} else {
				loggerServer.PrintInfof("client [%s] bind info [%s] error \n", c.GetKey(), message)
			}
		} else {
			if c.GetChecked() {
				param := make(map[string]interface{})
				err := json.Unmarshal([]byte(message), &param)
				if err != nil {
					loggerServer.PrintInfof("client [%s] parse param [%s] error\n", c.GetKey(), message)
					c.SendMsgf(net.CodeError, "code [%d] param [%s] error", code, message)
				} else {
					index := param[net.NameIndex]
					key := param[net.NameKey]
					callback := param[net.NameCallback]
					filepath := path.Join(conf.DbPath, fmt.Sprintf("%.f/%.f", index.(float64), key.(float64)))
					switch code {
					case net.CodeGet:
						loggerServer.PrintInfo(filepath)
						data, eGet := store.GetData(filepath)
						result := make(map[string]interface{})
						result["c"] = callback
						if eGet != nil {
							result["e"] = eGet
						} else {
							result["d"] = data
						}
						bResult, e := json.Marshal(result)
						if e != nil {
							logger.PrintError(e)
						}
						c.Send(net.CodeGet, bResult)
					case net.CodeSet:
						result := make(map[string]interface{})
						result["c"] = callback
						value := param[net.NameValue]
						eSave := store.Save(filepath, value.(map[string]interface{}))
						if eSave != nil {
							result["e"] = eSave
						}
						bResult, e := json.Marshal(result)
						if e != nil {
							logger.PrintError(e)
						}
						c.Send(net.CodeSet, bResult)
					case net.CodeDelete:
					}
				}
			} else {
				loggerServer.PrintInfof("client [%s] not binded\n", c.GetKey())
				c.SendMsg(net.CodeError, "not bind")
			}
		}
		// logger.PrintInfof("get client [%s] code [%d] message [%s] \n", c.GetKey(), code, message)
		// c.Send([]byte("from server get ["+message+"]"), -100)
	})
	server.OnClientClosed(func(c *net.Client) {
		loggerServer.PrintInfof("client [%s] closed\n", c.GetKey())
	})
	server.Listen(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
}
