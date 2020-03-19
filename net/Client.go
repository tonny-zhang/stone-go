package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Client tcp client
type Client struct {
	conn        net.Conn
	isChecked   bool
	onConnError func(err error)
	onConnect   func()
	onClose     func()
	onMessage   func(code int16, message string)
}

// OnConnError Called after connect server when error
func (c *Client) OnConnError(callback func(err error)) {
	c.onConnError = callback
}

// OnConnect Called after connect server
func (c *Client) OnConnect(onConnect func()) {
	c.onConnect = onConnect
}

// OnClose Called after connect closed
func (c *Client) OnClose(onClose func()) {
	c.onClose = onClose
}

// OnMessage Called after get new message
func (c *Client) OnMessage(onMessage func(code int16, message string)) {
	c.onMessage = onMessage
}

// Send text message to client
func (c *Client) Send(code int16, data []byte) error {
	packer := new(Packer)
	packer.Code = code
	packer.Length = int32(len(data))
	packer.Msg = data

	b, e := packer.PackToByte()
	if e != nil {
		return e
	}
	_, e1 := c.conn.Write(b)
	return e1
}

// SendMsg send string msg
func (c *Client) SendMsg(code int16, msg string) error {
	return c.Send(code, []byte(msg))
}

// SendMsgf send string msg by format
func (c *Client) SendMsgf(code int16, format string, argv ...interface{}) error {
	return c.SendMsg(code, fmt.Sprintf(format, argv...))
}

// GetKey get key of socket
func (c *Client) GetKey() string {
	conn := c.GetConn()
	if conn != nil {
		return conn.RemoteAddr().String()
	}
	return ""
}

// GetConn get net.Conn
func (c *Client) GetConn() net.Conn {
	return c.conn
}

// Conn client connect to server
func (c *Client) Conn(host string, port int) {
	conn, e := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		if c.onConnError != nil {
			c.onConnError(e)
		}
	}
	c.conn = conn
	c.work()
}

// Close close tcp
func (c *Client) Close() error {
	return c.conn.Close()
}

// SetChecked set checked
func (c *Client) SetChecked() {
	c.isChecked = true
}

// GetChecked get client isChecked
func (c *Client) GetChecked() bool {
	return c.isChecked
}

// 第2种处理socket拆包粘包问题
func (c *Client) work2() {
	if c.onConnect != nil {
		c.onConnect()
	}

	for {
		// 1. 先读取头部消息
		// 2. 从头部消息里解析出编码及数据长度
		// 3. 根据数据长度读取数据

		headData := make([]byte, 4+4+2)
		if _, err := io.ReadFull(c.conn, headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}
		spliter := make([]byte, 4)
		dataBuff := bytes.NewReader(headData)

		var packer Packer
		if err := binary.Read(dataBuff, binary.LittleEndian, &spliter); err != nil {
			fmt.Println("unpack1 ", err)
			break
		}
		if err := binary.Read(dataBuff, binary.LittleEndian, &packer.Code); err != nil {
			fmt.Println("unpack2 ", err)
			break
		}
		if err := binary.Read(dataBuff, binary.LittleEndian, &packer.Length); err != nil {
			fmt.Println("unpack3 ", err)
			break
		}

		if packer.Length > 0 {
			packer.Msg = make([]byte, packer.Length)
			if _, err := io.ReadFull(c.conn, packer.Msg); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}

		if c.onMessage != nil {
			c.onMessage(packer.Code, string(packer.Msg))
		}
	}
}

// Read client data from channel
func (c *Client) work() {
	if c.onConnect != nil {
		c.onConnect()
	}

	reader := bufio.NewReader(c.conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		lenData := len(data)
		if atEOF && lenData == 0 {
			return 0, nil, nil
		}
		if lenData > 10 {
			indexSplit := bytes.IndexAny(data, SPLITER)
			if indexSplit >= 0 {
				length := int32(0)
				binary.Read(bytes.NewReader(data[6:10]), binary.LittleEndian, &length)

				indexNext := indexSplit + 4 + 2 + 4 + int(length)
				if indexNext <= lenData {
					return indexNext, data[:indexNext], nil
				}
			}
		}

		if atEOF {
			return lenData, data, nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		scannedPack := new(Packer)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		// log.Println(scannedPack)
		if c.onMessage != nil {
			c.onMessage(scannedPack.Code, string(scannedPack.Msg))
		}
	}

	if c.onClose != nil {
		c.onClose()
	}
}

// NewClient get instance of Client
func NewClient() *Client {
	return new(Client)
}
