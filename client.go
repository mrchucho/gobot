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
			log.Printf("*** %s is leaving %s\n", msg.Args(1), msg.Where())
			quit <- true
		case "PRIVMSG":
			log.Printf("*** Heard %s say \"%s\" in %s\n", msg.Prefix, msg.Content(), msg.Where())
			self.bot.Handle(msg)
		case "QUIT":
			log.Printf("*** %s quit.\n", msg.Prefix)
		case "PART":
			log.Printf("*** %s left %s.\n", msg.Prefix, msg.Where())
		default:
			// TODO most stuff isn't implemented yet, so just ignore.
			log.Printf("*** Unhandled Command: %s.\n", msg.Command)
		}
	}
}
