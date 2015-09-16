package gobot

import (
	"log"
	"strconv"
)

// RFC 2812
const (
	RPL_ENDOFMOTD = 376
	ERR_NOMOTD    = 422
)

type Client struct {
	bot *Bot
}

func NewClient(bot *Bot) *Client {
	return &Client{bot}
}

func (self *Client) Process(msg *Message, quit chan bool) {
	log.Printf("<-- %v\n", msg)
	c, err := strconv.Atoi(msg.Command)
	if err == nil {
		switch c {
		case ERR_NOMOTD, RPL_ENDOFMOTD:
			self.bot.Join(self.bot.Channel)
		}
	} else {
		switch msg.Command {
		case "PING":
			self.bot.Pong(msg.Args(2))
		case "KICK":
			log.Printf("*** %s is leaving %s\n", msg.Args(1), msg.Args(0))
			quit <- true
		case "PRIVMSG":
			log.Printf("*** Heard %s say \"%s\" in %s\n", msg.Prefix, msg.Args(2), msg.Args(0))
			self.bot.Handle(msg)
		case "QUIT":
			log.Printf("*** %s quit.\n", msg.Args(2))
		case "PART":
			log.Printf("*** %s left %s.\n", msg.Args(2), msg.Args(0))
		default:
			// TODO most stuff isn't implemented yet, so just ignore.
			log.Printf("*** Invalid Command: %s.\n", msg.Command)
		}
	}
}

/* TODO - unused
func (self *Client) isForMe(msg *Message) (forMe bool, from string) {
	if msg.Args(0) == self.bot.Nick {
		forMe = true
		from = msg.Prefix
	} else if strings.HasPrefix(msg.Args(2), self.bot.Nick+":") {
		forMe = true
		from = msg.Args(0)
	} // else if HasPrefix("!")
	return
}
*/
