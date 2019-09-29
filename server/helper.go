package server

import (
	"encoding/base64"
	"log"
	"net"
)

func closingConnection(conn net.Conn, rmtAddr string) {
	log.Printf("CLIENT_DISCONNECT [%s] >> closing connection\n", rmtAddr)
	conn.Close()
}

func decodeBase64(str string) string {
	b, _ := base64.RawURLEncoding.DecodeString(str)
	return string(b)
}

func encodeBase64(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}
