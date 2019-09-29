package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
)

type initial struct {
	Username string
	Key      string
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

func StartBrokerService(port uint, key string) {
	listClient := ListClient{}
	tcp, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicln("port is used")
	}
	l, err := net.ListenTCP("tcp4", tcp)
	if err != nil {
		log.Panicln("cannot start the service")
	}
	fmt.Printf("-- broker server started (port: %d) --\n", port)
	for {
		conn, _ := l.Accept()
		rmtAddr := conn.RemoteAddr().String()
		log.Printf("CLIENT_CONNECT [%s] >> connected\n", rmtAddr)
		go func() {
			buffer := make([]byte, 1024*4)
			var init initial
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					closingConnection(conn, rmtAddr)
					return
				}
				b := buffer[0:n]
				decbuffer := decodeBase64(string(b))
				if err := json.Unmarshal([]byte(decbuffer), &init); err != nil {
					log.Printf("ERROR [%s] >> initial format wrong\n", rmtAddr)
				}
				if init.Username != "" && init.Key == key {
					break
				} else {
					log.Printf("CLIENT_VIOLATION [%s] >> wrong key or empty username\n", rmtAddr)
					closingConnection(conn, rmtAddr)
					return
				}
			}

			listClient.AddClient(Client{
				Username: init.Username,
				Conn:     &conn,
				IP:       rmtAddr,
			})
			log.Printf("CLIENT_JOIN [%s] >> as %s\n", rmtAddr, init.Username)

			for {
				arrDat := []string{}
				n, err := conn.Read(buffer)
				if err != nil {
					log.Printf("CLIENT_DISCONNECT [%s] >> closing connection\n", rmtAddr)
					log.Printf("CLIENT_REMOVE >> remove %s from list\n", rmtAddr)
					listClient.RemoveClient(rmtAddr)
					conn.Close()
					return
				}
				b := buffer[0:n]
				r, _ := regexp.Compile("(?:[A-z0-9-_]+)")
				arrDat = r.FindAllString(string(b), -1)
				var packet Packet
				for i := 0; i < len(arrDat); i++ {
					decData := decodeBase64(arrDat[i])
					if err := json.Unmarshal([]byte(decData), &packet); err != nil {
						log.Printf("ERROR [%s] >> packet format wrong\n", rmtAddr)
					}
					switch packet.Action {
					case "send":
						packet.Payload.From = init.Username
						log.Printf("CLIENT_SEND [%s] >> sending data to: %s, from: %s, data: %s", rmtAddr, packet.Payload.To, packet.Payload.From, packet.Payload.Data)
						listClient.SendData(packet.Payload)
					case "list":
						packet.Payload.To = init.Username
						ls := []string{}
						for _, v := range listClient {
							ls = append(ls, v.Username)
						}
						list, _ := json.Marshal(ls)
						packet.Payload.Data = string(list)
						log.Printf("CLIENT_SEND [%s] >> sending list to: %s", rmtAddr, init.Username)
						listClient.SendData(packet.Payload)
					}
				}
			}
		}()
	}
}
