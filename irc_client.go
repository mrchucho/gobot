package irc_client

import (
	"./irc";
	"./irc_bot";
	"log";
	"strconv";
	"strings";
)

// RFC 2812
const (
	RPL_ENDOFMOTD = 376;
)

type Client struct {
	bot *irc_bot.Bot;
}

func NewClient(bot *irc_bot.Bot) *Client {
	return &Client{bot};
}

func (self *Client) Process(msg *irc.Message, quit chan bool) {
	log.Stdoutf("<-- %v\n", msg);
	c, err := strconv.Atoi(msg.Command);
	if err == nil {
		switch c {
			case RPL_ENDOFMOTD:
				self.bot.Join(self.bot.Channel);
		}
	} else {
		switch msg.Command {
			case "PING":
				self.bot.Pong(msg.Args(2));
			case "KICK":
				log.Stdoutf("*** %s is leaving %s\n", msg.Args(1), msg.Args(0));
				quit <- true
			case "PRIVMSG":
				log.Stdoutf("*** Heard %s say \"%s\" in %s\n", msg.Prefix, msg.Args(2), msg.Args(0));
				self.bot.Handle(msg);
			case "QUIT":
				log.Stdoutf("*** %s quit.\n", msg.Args(2));
			case "PART":
				log.Stdoutf("*** %s left %s.\n", msg.Args(2), msg.Args(0));
			default:
				// TODO most stuff isn't implemented yet, so just ignore.
				// log.Stderrf("*** Invalid Command: %s.\n", msg.Command);
		}
	}
}

// FIXME - unused
func (self *Client) isForMe(msg *irc.Message) (forMe bool, from string) {
	if msg.Args(0) == self.bot.Nick {
		forMe = true;
		from  = msg.Prefix;
	} else if strings.HasPrefix(msg.Args(2), self.bot.Nick+":") {
		forMe = true;
		from  = msg.Args(0);
	} // else if HasPrefix("!")
	return;
}

/*
prefix - "who", e.g.
	:server.net - a message from the server
	:nick!~user@ip.address.net - message from a user

command - ### or IRC_COMMAND
params: usually of the format:
	<recipient> :<contents>
	#channel :blah blah blah
	user :blah blah blah
	#channel user : --- as in KICK
*/

