package gobot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

const (
	CarriageReturn = 0x0A
	LineFeed       = 0x0D
	Space          = 0x20
	Colon          = 0x3B
)

type Bot struct {
	Nick, User, Mode, RealName, Channel string
	Connection                          *net.Conn

	request  chan *Message
	response chan string
	handlers map[string]func(*Message, []string, *string)
}

func NewBot(nick, user, mode, realname, channel string, connection *net.Conn) *Bot {
	bot := &Bot{Nick: nick, User: user, Mode: mode, RealName: realname, Channel: channel, Connection: connection}
	bot.makeHandlerMap()
	return bot
}

func (bot *Bot) Run(client *Client) {
	quit := make(chan bool)
	bot.request = make(chan *Message)
	bot.response = make(chan string)
	reader := bufio.NewReader(*bot.Connection)

	bot.Login()

	for {
		go bot.getServerInput(reader, quit)
		select {
		case messageToServer := <-bot.response:
			go bot.write(messageToServer)
		case messageFromServer := <-bot.request:
			go client.Process(messageFromServer, quit)
		case <-quit:
			// shutdown properly
			return
		}
	}
}

func (self *Bot) getServerInput(reader *bufio.Reader, quit chan bool) {
	line, err := reader.ReadString(CarriageReturn)
	if err != nil {
		log.Printf("ERROR Reading: ", err)
		quit <- true
	} else {
		self.request <- NewMessage(line)
	}
}

func (self *Bot) send(command string) {
	self.response <- command
}

func (self *Bot) sendNow(command string) {
	self.write(command)
}

// FIXME enforce IRC 512 char. limit...
func (self *Bot) write(message string) {
	log.Printf("--> %s\n", message)
	fmt.Fprintf(*self.Connection, "%s\r\n", message)
}

// ------------------ IRC COMMANDS --------------------
func (self *Bot) Login() {
	self.sendNow(fmt.Sprintf("NICK %s", self.Nick))
	self.sendNow(fmt.Sprintf("USER %s %s %s %s", self.User, self.Mode, "*", self.RealName))
}

func (self *Bot) Join(channel string) {
	self.send(fmt.Sprintf("JOIN %s", channel))
}

func (self *Bot) Pong(pong string) {
	self.send(fmt.Sprintf("PONG %s", pong))
}

func (self *Bot) Say(what, where string) {
	self.send(fmt.Sprintf("PRIVMSG %s :%s", where, what))
}

// FIXME is this even right?
func (self *Bot) Quit(why string) {
	self.send(fmt.Sprintf("QUIT :%s", why))
}

// ----------------- "Command Handlers" ------------------------
func (self *Bot) makeHandlerMap() {
	self.handlers = map[string]func(*Message, []string, *string){
		"hello": func(msg *Message, args []string, where *string) {
			self.Say(fmt.Sprintf("Hi, %s!", msg.Prefix), *where)
		},
		"version": func(msg *Message, args []string, where *string) {
			self.Say("Version 0.0 Alpha", *where)
		},
		"join": func(msg *Message, args []string, where *string) {
			self.Join(args[0])
		},
		"quit": func(msg *Message, args []string, where *string) {
			self.Quit("Leaving because you asked.")
		},
		"sleep": func(msg *Message, args []string, where *string) {
			// prove you can multi-task, gobot!
			self.Say("Going to sleep...", *where)
			time.Sleep(10e9)
			self.Say("Awoke!", *where)
		},
	}
}

func (self *Bot) Handle(msg *Message) {
	if command, args, where := msg.GetCommand(&self.Nick); command != nil {
		if f, ok := self.handlers[*command]; ok {
			f(msg, args, where)
		}
	}
}