package irc

import (
		"net";
		"log";
		"bufio";
		"strings";
		"strconv";
		fmt "fmt";
)

const (
		CarriageReturn = 0x0A;
		LineFeed = 0x0D;
)

type Irc struct {
		Nick, User, Mode, RealName string;
		Connection net.Conn;
}

func (self *Irc) Run() {
	self.login();
	reader := bufio.NewReader(self.Connection);
	for {
		line, err := reader.ReadString(CarriageReturn);
		if err != nil {
			log.Stderr("ERROR Reading: ", err);
			break
		}
		self.handle(self.parse(line));
	}
}

func (self *Irc) login() {
	self.send(fmt.Sprintf("NICK %s", self.Nick));
	self.send(fmt.Sprintf("USER %s %s %s %s", self.User, self.Mode, "*", self.RealName));
}

func (self *Irc) handle(prefix string, command string, params string) {
	log.Stdoutf("<-- [%s][%s] %s\n", prefix, command, params);
	c, err := strconv.Atoi(command);
	if err == nil {
		switch c {
			case 376:
				log.Stdout("*** Greeting ended, join.");
				self.send("JOIN #test");
		}
	} else {
		log.Stdout("*** Not a numeric command: ", command);
		// this is where we'll handle commands
		switch command {
			case "PING":
				self.send(fmt.Sprintf("PONG %s", params));
			case "KICK":
				log.Stdout("*** Leaving");
			return
		}
	}
}

func (self *Irc) parse(msg string) (prefix string, command string, params string) {
	params = msg[0:len(msg)-1]; // chomp 
	if strings.HasPrefix(msg, ":") {
		prefix = params[0:strings.Index(params, " ")];
		params = params[len(prefix)+1:len(params)];
	}
	command = params[0:strings.Index(params, " ")];
	params = params[len(command)+1:len(params)];
	return
}

func (self *Irc) send(command string) {
	log.Stdoutf("--> %s\n", command);
	self.Connection.Write(strings.Bytes(command + "\r\n"));
}
