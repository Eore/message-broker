package main

import (
	"encoding/json"
	"fmt"
	"net"
)

//error code number
const (
	CodeNoError = iota
	CodeInvalidFormat
)

type Command struct {
	Action    string      `json:"action"`
	To        string      `json:"to,omitempty"`
	Parameter interface{} `json:"parameter"`
}

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data"`
}

type Client struct {
	UID  string `json:"uid"`
	IP   string `json:"ip"`
	conn *net.Conn
}

func Respond(conn net.Conn, res Response) {
	b, _ := json.Marshal(res)
	conn.Write(append(b, '\n'))
}

func main() {
	listClient := []Client{}
	l, _ := net.Listen("tcp4", ":9191")
	for {
		c, _ := l.Accept()
		fmt.Println(c.RemoteAddr().String())
		defer c.Close()
		go func(c net.Conn) {
			for {
				data := make([]byte, (1024 * 10))
				n, _ := c.Read(data)
				var cmd Command
				if err := json.Unmarshal(data[0:n], &cmd); err != nil {
					Respond(c, Response{
						Code: CodeInvalidFormat,
						Data: "format data salah",
					})
				} else {
					switch cmd.Action {
					case "join":
						fmt.Println("join")
						listClient = append(listClient, Client{
							UID:  fmt.Sprint(cmd.Parameter),
							IP:   c.RemoteAddr().String(),
							conn: &c,
						})
						Respond(c, Response{
							Code: CodeNoError,
							Data: cmd,
						})
					case "list":
						Respond(c, Response{
							Code: CodeNoError,
							Data: listClient,
						})
					case "send":
						for _, val := range listClient {
							if val.UID == cmd.To {
								Respond(*val.conn, Response{
									Code: CodeNoError,
									Data: cmd.Parameter,
								})
							}
						}

					default:
						fmt.Println("defa")
					}
				}
			}
		}(c)
	}
}
