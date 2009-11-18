package irc_bot

import (
		"./irc";
		"net";
		"log";
		"bufio";
		"strings";
		"fmt";
)

const (
		CarriageReturn = 0x0A;
		LineFeed = 0x0D;
)

type Bot struct {
		Nick, User, Mode, RealName, Channel string;
		Connection *net.Conn;
		request  chan *irc.Message;
		response chan string;
}

// accept os.Args
func NewBot(nick, user, mode, realname, channel string, connection *net.Conn) *Bot {
	return &Bot{Nick:nick, User:user, Mode:mode, RealName:realname, Channel:channel, Connection:connection};
}

func (bot *Bot) Run(client *irc.Client) {
	quit := make(chan bool);
	bot.request = make(chan *irc.Message);
	bot.response = make(chan string);
	reader := bufio.NewReader(*bot.Connection);

	bot.Login();

	for {
		go bot.getServerInput(reader, quit);
		select {
			case messageToServer := <-bot.response:
				go bot.write(messageToServer);
			case messageFromServer := <-bot.request:
				go client.Process(messageFromServer, quit);
			case <- quit:
				// shutdown properly
				return
		}
	}
}

func (self *Bot) getServerInput(reader *bufio.Reader, quit chan bool) {
	line, err := reader.ReadString(CarriageReturn);
	if err != nil {
		log.Stderr("ERROR Reading: ", err);
		quit <- true;
	}
	self.request <- self.parse(line);
}

func (self *Bot) parse(msg string) (ircMessage *irc.Message) {
	var prefix, command, params string;
	// EOF checking...
	params = msg[0:len(msg)-1]; // chomp 
	if strings.HasPrefix(msg, ":") {
		prefix = params[0:strings.Index(params, " ")];
		params = params[len(prefix)+1:len(params)];
	}
	command = params[0:strings.Index(params, " ")];
	params = params[len(command)+1:len(params)];
	return irc.NewMessage(prefix, command, params);
}

func (self *Bot) send(command string) {
	self.response <- command;
}

func (self *Bot) sendNow(command string) {
	self.write(command);
}

func (self *Bot) write(message string) {
	log.Stdoutf("--> %s\n", message);
	self.Connection.Write(strings.Bytes(message + "\r\n"));
}

// ------------------ IRC COMMANDS --------------------
func (self *Bot) Login() {
	self.sendNow(fmt.Sprintf("NICK %s", self.Nick));
	self.sendNow(fmt.Sprintf("USER %s %s %s %s", self.User, self.Mode, "*", self.RealName));
}

func (self *Bot) Join(channel string) {
	self.send(fmt.Sprintf("JOIN %s", channel));
}

func (self *Bot) Pong(pong string) {
	self.send(fmt.Sprintf("PONG %s", pong));
}
