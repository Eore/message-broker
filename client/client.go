package client

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net"
	"regexp"
	"time"
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

type Packet struct {
	Action    string
	Parameter string
	Payload   Payload
}

type Payload struct {
	UID   string `json:"uid"`
	To    string `json:"to"`
	From  string `json:"from"`
	Error string `json:"error"`
	Data  string `json:"data"`
}
type connection struct {
	username string
	Conn     net.Conn
}

type initial struct {
	Username string
	Key      string
}

func encodeBase64(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

func decodeBase64(str string) string {
	b, _ := base64.RawURLEncoding.DecodeString(str)
	return string(b)
}

func checkConnection() {}

func ClientConnect(host, username, key string) connection {
	log.Printf("connecting to %s\n", host)
	conn, err := net.Dial("tcp4", host)
	if err != nil {
		log.Fatalf("%s cant be reach\n", host)
		return connection{Conn: nil}
	}
	b, _ := json.Marshal(initial{
		Username: username,
		Key:      key,
	})
	str := encodeBase64(string(b))
	conn.Write([]byte(str))
	time.Sleep(time.Second)
	log.Printf("connected to %s\n", host)
	return connection{username: username, Conn: conn}
}

func (c *connection) SendData(uid, to, data string) {
	b, _ := json.Marshal(Packet{
		Action: "send",
		Payload: Payload{
			UID:  uid,
			To:   to,
			Data: data,
		},
	})
	str := encodeBase64(string(b))
	c.Conn.Write([]byte(str + "."))
}

func (c *connection) ListenData() []Payload {
	buffer := make([]byte, 1024*4)
	n, err := c.Conn.Read(buffer)
	if err == io.EOF || err != nil {
		log.Fatalln("disconnected from server")
		c.Conn.Close()
		return []Payload{}
	}
	b := buffer[0:n]
	r, _ := regexp.Compile("(?:[A-z0-9-_]+)")
	arr := r.FindAllString(string(b), -1)
	arrPayload := []Payload{}
	var payload Payload
	for i := range arr {
		json.Unmarshal([]byte(decodeBase64(arr[i])), &payload)
		arrPayload = append(arrPayload, payload)
	}
	return arrPayload
}

func (c *connection) SendDataWithResponse(uid, to, data string) Payload {
	c.SendData(uid, to, data)
	payload := c.ListenData()
	for i := range payload {
		if payload[i].UID == uid {
			return payload[0]
		}
	}
	return Payload{}
}
