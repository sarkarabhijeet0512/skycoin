package clients

import (
	"net"
)

type Client struct {
	Conn     net.Conn
	Name     string
	Room     *Room
	Commands chan<- Command
}

type Room struct {
	Name    string
	Members map[net.Addr]*Client
}
