package main

import (
	"encoding/base64"
	"flag"
	"log"
	"strings"

	irc "github.com/thoj/go-ircevent"

	"github.com/songgao/water"
)

func main() {
	channel := flag.String("channel", "#ipircdata", "")
	nick := flag.String("nick", "owo", "")
	addr := flag.String("connect", "", "")

	flag.Parse()

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("if: %s\n", ifce.Name())

	ircConn := irc.IRC(*nick, "this is cursed https://cynthia.re")

	ircConn.AddCallback("001", func(e *irc.Event) {
		ircConn.Join(*channel)
	})

	ircConn.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := e.Message()
		if strings.HasPrefix(msg, "PACKET ") {
			packet, err := base64.StdEncoding.DecodeString(msg[7:])
			if err != nil {
				log.Println(err)
			}
			log.Printf("Packet recv: % x\n", packet)
			_, err = ifce.Write(packet)
			if err != nil {
				log.Println("failed to write to if")
			}
		}
	})

	err = ircConn.Connect(*addr)
	if err != nil {
		log.Fatal(err)
	}

	go ircConn.Loop()

	/*go func(conn *irc.Connection) {
		packet := make([]byte, 2000)
		for {
			n, err := conn.Read(packet)
			if err != nil {
				log.Print("owo")
				//conn.Close()
				//break
			}

			log.Printf("Packet recv: % x\n", packet[:n])
			_, err = ifce.Write(packet[:n])
			if err != nil {
				log.Println("failed to write to if")
			}
		}
	}(ircConn)*/

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}

		ircConn.Privmsgf(*channel, "PACKET %s", base64.StdEncoding.EncodeToString(packet[:n]))
	}
}
