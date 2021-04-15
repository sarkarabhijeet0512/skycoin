package clients

import (
	"bufio"
	"fmt"
	"strings"
)

type CommandID int

const (
	CMD_NAME CommandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type Command struct {
	Id      CommandID
	Clients *Client
	Args    []string
}

// type client Client

func (c *Client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "/name":
			c.Commands <- Command{
				Id:      CMD_NAME,
				Clients: c,
				Args:    args,
			}
		case "/join":
			c.Commands <- Command{
				Id:      CMD_JOIN,
				Clients: c,
				Args:    args,
			}
		case "/rooms":
			c.Commands <- Command{
				Id:      CMD_ROOMS,
				Clients: c,
				Args:    args,
			}
		case "/msg":
			c.Commands <- Command{
				Id:      CMD_MSG,
				Clients: c,
				Args:    args,
			}
		case "/quit":
			c.Commands <- Command{
				Id:      CMD_QUIT,
				Clients: c,
				Args:    args,
			}
		default:
			c.Err(fmt.Errorf("unknown Command %s", cmd))
		}

	}
}
func (c *Client) Err(err error) {
	c.Conn.Write([]byte("ERR:" + err.Error() + "\n"))
}
func (c *Client) Msg(msg string) {
	c.Conn.Write([]byte(">" + msg + "\n"))
}
func (r *Room) Broadcast(sender *Client, msg string) {
	for addr, m := range r.Members {
		if addr != sender.Conn.RemoteAddr() {
			m.Msg(msg)
		}
	}
}
