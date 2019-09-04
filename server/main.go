package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		rmtAddr := c.RemoteAddr().String()
		fmt.Println(rmtAddr)
		// defer c.Close()
		go func(c net.Conn) {
			for {
				data := make([]byte, (1024 * 10))
				n, err := c.Read(data)
				if err == io.EOF {
					fmt.Println("closing", rmtAddr)
					for i, val := range listClient {
						if val.IP == rmtAddr {
							fmt.Println("remove", rmtAddr)
							listClient = append(listClient[:i], listClient[i+1:]...)
						}
					}
					c.Close()
				}
				var cmd Command
				if err := json.Unmarshal(data[0:n], &cmd); err != nil {
					Respond(c, Response{
						Code: CodeInvalidFormat,
						Data: "format data salah",
					})
				} else {
					switch cmd.Action {
					case "join":
						log.Println(rmtAddr, "join as", cmd.Parameter)
						listClient = append(listClient, Client{
							UID:  fmt.Sprint(cmd.Parameter),
							IP:   rmtAddr,
							conn: &c,
						})
						Respond(c, Response{
							Code: CodeNoError,
							Data: cmd,
						})
					case "list":
						log.Println(rmtAddr, "calling list")
						Respond(c, Response{
							Code: CodeNoError,
							Data: listClient,
						})
					case "send":
						for _, val := range listClient {
							if val.UID == cmd.To {
								log.Println(rmtAddr, "sending to", val.IP, val.UID)
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
