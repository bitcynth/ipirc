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

	// setup the tun interface
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// just log the interface name
	log.Printf("if: %s\n", ifce.Name())

	// setup the irc connection
	ircConn := irc.IRC(*nick, "this is cursed https://cynthia.re")

	// callback for irc connection being established
	ircConn.AddCallback("001", func(e *irc.Event) {
		ircConn.Join(*channel)
	})

	// callback for messages on irc
	ircConn.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := e.Message()
		// do magic with packet messages
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

	// connect to irc
	err = ircConn.Connect(*addr)
	if err != nil {
		log.Fatal(err)
	}

	go ircConn.Loop()

	// listen for packets to send to irc
	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}

		ircConn.Privmsgf(*channel, "PACKET %s", base64.StdEncoding.EncodeToString(packet[:n]))
	}
}
