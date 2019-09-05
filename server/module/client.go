package module

import (
	"encoding/json"
	"errors"
	"net"
)

type Command struct {
	Action    string      `json:"action"`
	To        string      `json:"to,omitempty"`
	Parameter interface{} `json:"parameter"`
}

type Response struct {
	Code uint        `json:"code"`
	From string      `json:"from,omitempty"`
	Data interface{} `json:"data"`
}

type Client struct {
	UID  string    `json:"uid"`
	IP   string    `json:"ip"`
	Conn *net.Conn `json:"-"`
}

type ListClient []Client

func (l *ListClient) SendData(from string, cmd Command) {
	for _, client := range *l {
		if client.UID == cmd.To {
			b, _ := json.Marshal(Response{
				Code: 0,
				From: from,
				Data: cmd.Parameter,
			})
			(*client.Conn).Write(append(b, '\n'))
		}
	}
}

func (l *ListClient) CariClient(uid string) Client {
	for _, client := range *l {
		if client.UID == uid {
			return client
		}
	}
	return Client{}
}

func (l *ListClient) TambahClient(c Client) error {
	for _, client := range *l {
		if client.UID == c.UID || client.IP == c.IP {
			return errors.New("client already exist")
		}
	}
	*l = append(*l, c)
	return nil
}

func (l *ListClient) HapusClient(ip string) error {
	for i, client := range *l {
		if client.IP == ip {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}
	return errors.New("client not found")
}
