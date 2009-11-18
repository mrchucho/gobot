package irc

import (
	"fmt";
	"strings";
)

type Message struct {
	Prefix, Command, Params string;
	args []string;
}

func NewMessage(prefix, command, params string) *Message {
	// maybe parse Args here
	// set from, to, channel
	return &(Message{Prefix: prefix, Command: command, Params:params});
}

func (self *Message) String() string {
	return fmt.Sprintf("[%s][%s] %s", self.Prefix, self.Command, self.Params);
}

func (self *Message) Args(index int) string {
	// needs error handling
	if self.args == nil {
		// self.args = new([3]string);
		self.args = make([]string, 3);
		colonAt := strings.Index(self.Params, ":");
		for i, a := range(strings.Split(self.Params[0:colonAt], " ", 0)) {
			self.args[i] = a;
		}
		end := len(self.args);
		self.args[end-1] = self.Params[colonAt+1:len(self.Params)-1];
	}
	return self.args[index];
}

type Client interface {
	Process(*Message, chan bool);
}

