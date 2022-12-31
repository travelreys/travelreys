package utils

import "github.com/nats-io/nats.go"

func MakeNATSConn(url string) (*nats.Conn, error) {
	return nats.Connect(url, nats.Name("TiinyPlanet"))
}
