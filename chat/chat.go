package main

import (
	"bufio"
	"github.com/mattn/go-pubsub"
	"log"
	"net"
	"strings"
)

func main() {
	ps := pubsub.New()

	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listing", l.Addr().String())
	log.Println("Clients can connect to this server like follow:")
	log.Println("  $ telnet server:5555")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			buf := bufio.NewReader(c)
			ps.Sub(func(t string) {
				log.Println(t)
				c.Write([]byte(t + "\n"))
			})
			log.Println("Subscribed", c.RemoteAddr().String())
			for {
				b, _, err := buf.ReadLine()
				if err != nil {
					log.Println("Closed", c.RemoteAddr().String())
					break
				}
				ps.Pub(strings.TrimSpace(string(b)))
			}
			ps.Leave(nil)
		}(c)
	}
}
