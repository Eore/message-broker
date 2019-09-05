package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/eore/net-tcp/server/module"
)

//error code number
const (
	CodeNoError = iota
	CodeInvalidFormat
	CodeClientExist
)

func Respond(conn net.Conn, res module.Response) {
	b, _ := json.Marshal(res)
	conn.Write(append(b, '\n'))
}

func main() {
	var listClient module.ListClient
	l, _ := net.Listen("tcp4", ":9191")
	for {
		c, _ := l.Accept()
		rmtAddr := c.RemoteAddr().String()
		log.Printf(">> %s connected", rmtAddr)
		go func(c net.Conn) {
			var user string
			for {
				data := make([]byte, (1024 * 10))
				n, err := c.Read(data)
				if err == io.EOF {
					log.Printf(">> %s closing connection", rmtAddr)
					log.Printf(">> remove %s from list", rmtAddr)
					listClient.HapusClient(rmtAddr)
					c.Close()
				}
				var cmd module.Command
				if err := json.Unmarshal(data[0:n], &cmd); err != nil {
					Respond(c, module.Response{
						Code: CodeInvalidFormat,
						From: "server",
						Data: "format data salah",
					})
				} else {
					switch cmd.Action {
					case "join":
						log.Printf(">> %s join as %s", rmtAddr, cmd.Parameter)
						user = fmt.Sprint(cmd.Parameter)
						if err := listClient.TambahClient(module.Client{
							UID:  fmt.Sprint(cmd.Parameter),
							IP:   rmtAddr,
							Conn: &c,
						}); err != nil {
							Respond(c, module.Response{
								Code: CodeClientExist,
								From: "server",
								Data: err.Error(),
							})
						} else {
							Respond(c, module.Response{
								Code: CodeNoError,
								From: "server",
								Data: fmt.Sprintf("joined as %s", cmd.Parameter),
							})
						}
					case "list":
						log.Printf(">> %s calling list", rmtAddr)
						Respond(c, module.Response{
							Code: CodeNoError,
							From: "server",
							Data: listClient,
						})
					case "send":
						listClient.SendData(user, cmd)
					default:
						fmt.Println("defa")
					}
				}
			}
		}(c)
	}
}
