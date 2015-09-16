package gobot

import (
	"fmt"
	"regexp"
	"strings"
)

// http://tools.ietf.org/html/rfc2812#section-2.3
/*
Prefix - "who", optional. e.g.
	:server.net - a message from the server
	:nick!~user@ip.address.net - message from a user
Command - ### or IRC_COMMAND
Params: usually of the format:
	<recipient> :<contents>
	#channel :blah blah blah
	user :blah blah blah
	#channel user : --- as in KICK
*/
type Message struct {
	Prefix, Command, Params string
	args                    []string // Param string args.
	re                      *regexp.Regexp
}

// Parse the string: [Prefix (OPTIONAL)][Command][Parameters] and remove \r\n
func NewMessage(msg string) *Message {
	noticeRe := regexp.MustCompile(`^(NOTICE|ERROR) (.*)$`)
	if parsedMsg := noticeRe.FindAllString(msg, -1); len(parsedMsg) == 3 {
		return &(Message{
			Prefix: "",
			Command: parsedMsg[1],
			Params: parsedMsg[2][0:len(parsedMsg[2])-2]})
	} else {
		parsedMsg := strings.SplitN(msg, " ", 3)
		if len(parsedMsg) == 3 {
			return &(Message{
				Prefix: parsedMsg[0][1:len(parsedMsg[0])],
				Command: parsedMsg[1],
				Params: parsedMsg[2][0:len(parsedMsg[2])-2]})
		} else {
			return &(Message{
				Prefix: "", // No Prefix
				Command: parsedMsg[0],
				Params: parsedMsg[1][0:len(parsedMsg[1])-1]})
		}
	}
}

func (self *Message) String() string {
	return fmt.Sprintf("[%s][%s][%s]", self.Prefix, self.Command, self.Params)
}

// FIXME
func (self *Message) Args(index int) string {
	// needs error handling
	if self.args == nil {
		self.args = make([]string, 3)
		colonAt := strings.Index(self.Params, ":")
		for i, a := range strings.Split(self.Params[0:colonAt], " ") {
			self.args[i] = a
		}
		end := len(self.args)
		self.args[end-1] = self.Params[colonAt+1 : len(self.Params)]
	}
	return self.args[index]
}

// FIXME and/or merge w/ Args... also return :from
// FIXME "Command" is a bad name... this isn't the Command, it's the
// "instruction" to nick.
/*
#test :gobot hello
*/
// Maybe create another class... instruction? or do we even care if it's to
// nick?
// this should return the thing w/ command/args/where OR decorate that part of
// the Message
func (self *Message) GetCommand(nick *string) (command *string, args []string) {
	if self.re == nil {
		//						      1      2       3       4
		//						      where :nick    command args
		self.re = regexp.MustCompile(`^(.*) :(\w+)\s*(\w+)\s*(.*)$`)
	}
	m := self.re.FindStringSubmatch(self.Params)
	if len(m) >= 3 && m[2] == *nick {
		command = &m[3]
		args = strings.Split(m[4], " ")
	}
	return
}

func (self *Message) Where() string {
	return self.Args(0)
}

func (self *Message) Content() string {
	return self.Args(2)
}
