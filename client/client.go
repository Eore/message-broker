package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type Data struct {
	ID     string      `json:"id"`
	Error  Error       `json:"error"`
	Type   string      `json:"type"`
	Method string      `json:"method"`
	From   string      `json:"from"`
	To     string      `json:"to"`
	Body   interface{} `json:"body"`
}

type connection struct {
	Conn net.Conn
}

func ClientConnect(host string) connection {
	fmt.Printf("connecting to %s\n", host)
	conn, err := net.Dial("tcp4", host)
	if err != nil {
		fmt.Printf("%s cant be reach\n", host)
		return connection{Conn: nil}
	}
	fmt.Printf("connected to %s\n", host)
	return connection{Conn: conn}
}

func (c *connection) SendCommand(command string) {
	c.Conn.Write([]byte(command))
}

func (c *connection) SendData(data Data) {
	b, _ := json.Marshal(data)
	dat := fmt.Sprintf("send %s\n", b)
	c.Conn.Write([]byte(dat))
}

func (c *connection) ListenData(channel chan<- interface{}) {
	go func() {
		for {
			buffer := make([]byte, 1024*4)
			n, err := c.Conn.Read(buffer)
			if err == io.EOF {
				log.Fatalln("disconnected from server")
				c.Conn.Close()
				return
			}
			dataStr := strings.TrimSpace(string(buffer[0:n]))
			channel <- dataStr
		}
	}()
}

func (c *connection) ListenDataV2() string {
	buffer := make([]byte, 1024*4)
	n, err := c.Conn.Read(buffer)
	if err == io.EOF {
		log.Fatalln("disconnected from server")
		c.Conn.Close()
		return ""
	}
	return strings.TrimSpace(string(buffer[0:n]))
}
