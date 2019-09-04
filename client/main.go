package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Command struct {
	Action    string      `json:"action"`
	To        string      `json:"to,omitempty"`
	Parameter interface{} `json:"parameter"`
}

func main() {
	c, _ := net.Dial("tcp4", "localhost:9191")
	b, _ := json.Marshal(Command{
		Action:    "join",
		Parameter: "bot",
	})
	c.Write(b)
	time.Sleep(time.Second)
	bS, _ := json.Marshal(Command{
		Action:    "send",
		To:        "miaw",
		Parameter: "ini pesan",
	})
	c.Write(bS)
	// go func() {
	for {
		data := make([]byte, 1024*10)
		n, _ := c.Read(data)
		fmt.Println(string(data[0:n]))
	}
	// }()
	// fmt.Println(string(b))
	// fmt.Fprintln(c, string(b))
}
