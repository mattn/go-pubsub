package main

import (
	"fmt"
	"github.com/mattn/go-pubsub"
	"time"
)

type foo struct {
	bar string
}

func main() {
	ps := pubsub.New()
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

	ps.Pub(1)
	ps.Pub("hello")
	ps.Pub(2)
	ps.Pub(&foo{"bar!"})
	ps.Pub(&foo{"bar!"})

	time.Sleep(5 * time.Second)
}
