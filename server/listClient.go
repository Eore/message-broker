package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
)

func encodeBase64(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

type Client struct {
	Username string    `json:"username"`
	IP       string    `json:"ip"`
	Conn     *net.Conn `json:"-"`
}

type ListClient []Client

func (l *ListClient) SendData(data Payload) {
	for _, client := range *l {
		if client.Username == data.To {
			b, _ := json.Marshal(data)
			b64 := encodeBase64(string(b))
			(*client.Conn).Write(append([]byte(b64), '.'))
		}
	}
}

func (l *ListClient) FindClient(username string) Client {
	for _, client := range *l {
		if client.Username == username {
			return client
		}
	}
	return Client{}
}

func (l *ListClient) AddClient(c Client) error {
	for _, client := range *l {
		if client.Username == c.Username || client.IP == c.IP {
			return errors.New("client already exist")
		}
	}
	*l = append(*l, c)
	return nil
}

func (l *ListClient) RemoveClient(ip string) error {
	for i, client := range *l {
		if client.IP == ip {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}
	return errors.New("client not found")
}
