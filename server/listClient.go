package server

import (
	"encoding/json"
	"errors"
	"net"
)

type Client struct {
	Username string    `json:"username"`
	IP       string    `json:"ip"`
	Conn     *net.Conn `json:"-"`
}

type ListClient []Client

func (l *ListClient) SendData(data Data) {
	for _, client := range *l {
		if client.Username == data.To {
			b, _ := json.MarshalIndent(Data{
				ID: data.ID,
				Error: Error{
					Code:    ErrorNull,
					Message: "",
				},
				Type: data.Type,
				From: data.From,
				To:   data.To,
				Body: data.Body,
			}, "", "    ")
			(*client.Conn).Write(append(b, '\n'))
		}
	}
}

func (l *ListClient) CariClient(username string) Client {
	for _, client := range *l {
		if client.Username == username {
			return client
		}
	}
	return Client{}
}

func (l *ListClient) TambahClient(c Client) error {
	for _, client := range *l {
		if client.Username == c.Username || client.IP == c.IP {
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
