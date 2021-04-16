package server

import "sky_coin_messanger/clients"

type Server struct {
	Room    map[string]*clients.Room
	Command chan clients.Command
}
