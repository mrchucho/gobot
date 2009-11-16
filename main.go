package main

import . "./irc"
import (
	"net";
	"log";
)

func main() {
	conn, err := net.Dial("tcp", "", "67.207.138.175:6667");
	if err != nil {
		log.Exit("ERROR Dialing: ", err);
	}
	irc := Irc{Nick:"gobot", User:"rchurchil", Mode:"0", RealName:"Go Bot", Connection:conn};
	irc.Run();
	log.Stdoutf("Closing Connection...");
	conn.Close();
}
