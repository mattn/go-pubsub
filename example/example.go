package main

import (
	"fmt"
	"github.com/mattn/go-pubsub"
	"log"
	"time"
)

type foo struct {
	bar string
}

func main() {
	ps := pubsub.New()

	go func() {
		for err := range ps.Error() {
			if pse, ok := err.(*pubsub.PubSubError); ok {
				log.Println(pse.String())
				ps.Leave(pse.Subscriber())
			} else {
				log.Println(err.Error())
			}
		}
	}()

	ps.Sub(func(i int) {
		fmt.Println("int subscriber: ", i)
	})
	ps.Sub(func(s string) {
		fmt.Println("string subscriber: ", s)
	})
	ps.Sub(func(f *foo) {
		fmt.Println("foo subscriber1: ", f.bar)
	})
	ps.Sub(func(f *foo) {
		fmt.Println("foo subscriber2: ", f.bar)
		ps.Leave(nil)
	})
	var f3 func(f *foo)
	f3 = func(f *foo) {
		fmt.Println("foo subscriber3: ", f.bar)
		ps.Leave(f3)
	}
	ps.Sub(f3)
	ps.Sub(func(f float64) {
		fmt.Println("float64 subscriber: ", f)
		panic("Crash!")
	})

	ps.Pub(1)
	ps.Pub("hello")
	ps.Pub(2)
	ps.Pub(2.4)
	ps.Pub(&foo{"bar!"})
	time.Sleep(1 * time.Second)
	ps.Pub(2.4)
	ps.Pub(&foo{"bar!"})

	time.Sleep(5 * time.Second)
}
