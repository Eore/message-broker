package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Data struct {
	ID    string      `json:"id"`
	Error Error       `json:"error"`
	Type  string      `json:"type"`
	From  string      `json:"from"`
	To    string      `json:"to"`
	Body  interface{} `json:"body"`
}

func errorChecking(err error, fn func()) {
	if err != nil {
		log.Println(err)
	}
	if fn != nil {
		fn()
	}
}

func sendResponse(conn net.Conn, data Data) {
	data.Type = "response"
	b, _ := json.MarshalIndent(data, "", "    ")
	conn.Write(append(b, '\n'))
}

func StartBrokerServer(port string) {
	listClient := ListClient{}
	tcp, err := net.ResolveTCPAddr("tcp4", port)
	errorChecking(err, nil)
	l, err := net.ListenTCP("tcp4", tcp)
	errorChecking(err, nil)
	for {
		conn, _ := l.Accept()
		rmtAddr := conn.RemoteAddr().String()
		log.Printf(">> %s connected\n", rmtAddr)
		go func() {
			var username string
			for {
				buffer := make([]byte, 1024*4)
				n, err := conn.Read(buffer)
				if err != nil {
					log.Printf(">> %s closing connection", rmtAddr)
					log.Printf(">> remove %s from list", rmtAddr)
					listClient.HapusClient(rmtAddr)
					conn.Close()
					return
				}
				arrStr := strings.Split(strings.TrimSpace(string(buffer[0:n])), " ")
				switch arrStr[0] {
				case "send":
					if username == "" {
						sendResponse(conn, Data{
							ID:    Hashing(fmt.Sprint(time.Now())),
							Error: Error.New(Error{}, ErrorJoinFirst),
							From:  "server",
							To:    username,
							Body:  nil,
						})
					} else {
						d := Data{}
						json.Unmarshal([]byte(arrStr[1]), &d)
						d.Error = Error.New(Error{}, ErrorNull)
						d.From = username
						log.Printf(">> %s sending data to %s", d.From, d.To)
						listClient.SendData(d)
					}

				case "list":
					b, _ := json.MarshalIndent(Data{
						ID:    Hashing(fmt.Sprint(time.Now())),
						Error: Error.New(Error{}, ErrorNull),
						Type:  "response",
						From:  "server",
						To:    username,
						Body:  listClient,
					}, "", "    ")
					conn.Write(append(b, '\n'))

				case "join":
					username = arrStr[1]
					err := listClient.TambahClient(Client{
						Username: username,
						IP:       rmtAddr,
						Conn:     &conn,
					})
					if err != nil {
						b, _ := json.MarshalIndent(Data{
							ID:    Hashing(fmt.Sprint(time.Now())),
							Error: Error.New(Error{}, ErrorClientExist),
							Type:  "response",
							From:  "server",
							To:    username,
							Body:  nil,
						}, "", "    ")
						conn.Write(append(b, '\n'))
					} else {
						log.Printf(">> %s join as %s", rmtAddr, username)
						b, _ := json.MarshalIndent(Data{
							ID:    Hashing(fmt.Sprint(time.Now())),
							Error: Error.New(Error{}, ErrorNull),
							Type:  "response",
							From:  "server",
							To:    username,
							Body:  fmt.Sprintf("joined as %s", username),
						}, "", "    ")
						conn.Write(append(b, '\n'))
					}
				}
			}
		}()
	}
}
