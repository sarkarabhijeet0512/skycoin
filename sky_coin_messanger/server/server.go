package server

import (
	"fmt"
	"log"
	"net"
	"sky_coin_messanger/clients"
	"strings"
)

func NewServer() *Server {
	return &Server{
		Room:    make(map[string]*clients.Room),
		Command: make(chan clients.Command),
	}
}
func (s *Server) Run() {
	for cmd := range s.Command {
		switch cmd.Id {
		case clients.CMD_NAME:
			s.Name(cmd.Clients, cmd.Args)
		case clients.CMD_JOIN:
			s.Join(cmd.Clients, cmd.Args)
		case clients.CMD_ROOMS:
			s.RoomList(cmd.Clients, cmd.Args)
		case clients.CMD_MSG:
			s.Msg(cmd.Clients, cmd.Args)
		case clients.CMD_QUIT:
			s.Quit(cmd.Clients, cmd.Args)
		}
	}
}
func (s *Server) NewClient(conn net.Conn) {
	log.Println("new client has connected", conn.RemoteAddr().String())

	c := &clients.Client{
		Conn:     conn,
		Name:     "anonymous",
		Commands: s.Command,
	}
	c.ReadInput()
}
func (s *Server) Name(c *clients.Client, args []string) {
	if len(args) < 2 {
		c.Msg("nick is required. usage: /name NAME")
		return
	}
	c.Name = args[1]
	c.Msg(fmt.Sprintf("all right, I will call you %s", c.Name))
}
func (s *Server) Join(c *clients.Client, args []string) {
	if len(args) < 2 {
		c.Msg("room name is required. usage: /join ROOM_NAME")
		return
	}
	roomName := args[1]
	r, ok := s.Room[roomName]
	if !ok {
		r = &clients.Room{
			Name:    roomName,
			Members: make(map[net.Addr]*clients.Client),
		}
		s.Room[roomName] = r
	}
	r.Members[c.Conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.Room = r
	r.Broadcast(c, fmt.Sprintf("%s joined the room", c.Name))
	c.Msg(fmt.Sprintf("welcome to %s", roomName))

}
func (s *Server) RoomList(c *clients.Client, args []string) {
	var rooms []string
	for name := range s.Room {

		rooms = append(rooms, name)
	}
	c.Msg(fmt.Sprintln("available rooms are", strings.Join(rooms, ",")))
}
func (s *Server) Msg(c *clients.Client, args []string) {
	if len(args) < 2 {
		c.Msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.Room.Broadcast(c, c.Name+": "+msg)
}
func (s *Server) Quit(c *clients.Client, args []string) {
	log.Println("client has left the room or disconnected", c.Conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.Msg("sad to see you go")
	c.Conn.Close()
}
func (s *Server) quitCurrentRoom(c *clients.Client) {
	if c.Room != nil {
		oldRoom := s.Room[c.Room.Name]
		delete(s.Room[c.Room.Name].Members, c.Conn.RemoteAddr())
		oldRoom.Broadcast(c, fmt.Sprintf("%s has left the room", c.Name))
	}
}
