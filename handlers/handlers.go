package handlers

import (
	"github.com/mrchucho/gobot"
	"regexp"
	"net/http"
	"io/ioutil"
	"html"
	"fmt"
)

type Links struct {
	gobot.BotHandler
	reTitle *regexp.Regexp
}

func NewLinks(bot *gobot.Bot) *Links {
	links := &Links{}
	links.Bot = bot
	links.Matcher = regexp.MustCompile(`^.*(https?://\S+).*$`)
	links.reTitle = regexp.MustCompile(`<\s*?title\s*?>(.*)<\s*?\/title\s*?>`)
	return links
}

func (self *Links) Handle(msg *gobot.Message) bool {
	if url, ok := self.Matchs(msg) ; ok {
		resp, err := http.Get(url[1])
		if err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body) ; err == nil {
				link := fmt.Sprintf(
					"%s -=[ %s ]=-",
					html.EscapeString(string(self.reTitle.FindSubmatch(body)[1])),
					url[1])
				self.Bot.Say(link, msg.Where())
			}
		}
		return true
	}
	return false
}

func Load(bot *gobot.Bot) bool {
	bot.Handlers = []gobot.Handler{
		NewLinks(bot),
	}
	return true
}
