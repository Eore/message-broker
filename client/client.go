package client

import (
	"encoding/json"
	"net"
	"strings"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type Data struct {
	ID    string      `json:"id"`
	Error Error       `json:"error"`
	Type  string      `json:"type"`
	From  string      `json:"from"`
	To    string      `json:"to"`
	Body  interface{} `json:"body"`
}

type connection struct {
	Conn net.Conn
}

func ClientConnect(host string) connection {
	conn, _ := net.Dial("tcp4", host)
	return connection{Conn: conn}
}

func (c *connection) SendCommand(command string) {
	c.Conn.Write([]byte(command))
}

func (c *connection) SendData(data Data) {
	b, _ := json.MarshalIndent(data, "", "    ")
	c.Conn.Write(append(b, '\n'))
}

func (c *connection) ListenData(channel chan<- interface{}) {
	go func() {
		for {
			buffer := make([]byte, 1024*4)
			n, err := c.Conn.Read(buffer)
			if err != nil {
				c.Conn.Close()
			}
			dataStr := strings.TrimSpace(string(buffer[0:n]))
			channel <- dataStr
		}
	}()
}

// type Command struct {
// 	Action    string      `json:"action"`
// 	To        string      `json:"to,omitempty"`
// 	Parameter interface{} `json:"parameter"`
// }

// func main() {
// 	c, _ := net.Dial("tcp4", "localhost:9191")
// 	b, _ := json.Marshal(Command{
// 		Action:    "join",
// 		Parameter: "bot",
// 	})
// 	c.Write(b)
// 	time.Sleep(time.Second)
// 	bS, _ := json.Marshal(Command{
// 		Action:    "send",
// 		To:        "miaw",
// 		Parameter: "ini pesan",
// 	})
// 	c.Write(bS)
// 	// go func() {
// 	for {
// 		data := make([]byte, 1024*10)
// 		n, _ := c.Read(data)
// 		fmt.Println(string(data[0:n]))
// 	}
// 	// }()
// 	// fmt.Println(string(b))
// 	// fmt.Fprintln(c, string(b))
// }
