package handlers

import (
	"fmt"
	"github.com/mrchucho/gobot"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Links struct {
	*gobot.BotHandler
	reTitle *regexp.Regexp
}

func NewLinks(bot *gobot.Bot) *Links {
	return &Links{
		BotHandler: &gobot.BotHandler{
			Bot:     bot,
			Matcher: regexp.MustCompile(`^.*(https?://\S+).*$`),
		},
		reTitle: regexp.MustCompile(`<\s*?title\s*?>(.*)<\s*?\/title\s*?>`),
	}
}

func (self *Links) Handle(msg *gobot.Message) bool {
	url, ok := self.Matchs(msg)
	if !ok {
		return false
	}
	resp, err := http.Get(url[1])
	if err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			link := fmt.Sprintf(
				"%s -=[ %s ]=-",
				html.EscapeString(string(self.reTitle.FindSubmatch(body)[1])),
				url[1])
			self.Bot.Say(link, msg.Where())
			return true
		}
	}
	return false
}

func Load(bot *gobot.Bot) bool {
	bot.Handlers = []gobot.Handler{
		NewLinks(bot),
	}
	return true
}
