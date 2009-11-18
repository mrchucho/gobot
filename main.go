package main

import (
	"./irc";
	"./irc_bot";
	"./irc_client";
	"net";
	"log";
	"flag";
	"fmt";
)

// TODO handle signals
func main() {
	var nick, user, name, host, port, channel, mode string;
	flag.StringVar(&nick, "nick", "gobot", "IRC nick");
	flag.StringVar(&user, "user", "", "Username");
	flag.StringVar(&name, "name", "Go Bot", "Realname");
	flag.StringVar(&host, "host", "", "IRC host");
	flag.StringVar(&port, "port", "6667", "Port");
	flag.StringVar(&channel, "channel", "", "IRC channel");
	flag.StringVar(&mode, "mode", "0", "IRC mode");
	flag.Parse();

	conn, err := net.Dial("tcp", "", fmt.Sprintf("%s:%s", host, port));
	if err != nil {
		log.Exit("ERROR Dialing: ", err);
	}
	bot := irc_bot.NewBot(nick, user, mode, name, channel, &conn); // pass in as Args
	client := irc.Client(irc_client.NewClient(bot));
	bot.Run(&client);
	log.Stdoutf("Closing Connection...");
	conn.Close();
}
