package gobot

import (
	"fmt"
	"regexp"
	"strings"
)

// http://tools.ietf.org/html/rfc2812#section-2.3
type Message struct {
	Prefix, Command, Params string
	args                    []string // Param string args.
	re                      *regexp.Regexp
}

func NewMessage(prefix, command, params string) *Message {
	// maybe parse Args here
	// set from, to, channel
	return &(Message{Prefix: prefix, Command: command, Params: params})
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
func (self *Message) GetCommand(nick *string) (command *string, args []string, where *string) {
	if self.re == nil {
		self.re = regexp.MustCompile(`^(.*) :((.*):?) (.*)$`)
	}
	m := self.re.FindStringSubmatch(self.Params)
	if len(m) >= 4 && m[3] == *nick {
		command_and_args := strings.Split(m[4], " ")
		command = &command_and_args[0]
		args = command_and_args[1:len(command_and_args)]
		where = &m[1]
	}
	return
}
