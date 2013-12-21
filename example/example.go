package main

import (
	"fmt"
	"github.com/mattn/go-multichannel"
	"time"
)

type foo struct {
	bar string
}

func main() {
	mc := multichannel.New()
	mc.Sub(func(i int) {
		fmt.Println("int subscriber: ", i)
	})
	mc.Sub(func(s string) {
		fmt.Println("string subscriber: ", s)
	})
	mc.Sub(func(f *foo) {
		fmt.Println("foo subscriber1: ", f.bar)
	})
	var f2 func(f *foo)
	f2 = func(f *foo) {
		fmt.Println("foo subscriber2: ", f.bar)
		mc.Leave(f2)
	}
	mc.Sub(f2)

	mc.Pub(1)
	mc.Pub("hello")
	mc.Pub(2)
	mc.Pub(&foo{"bar!"})
	mc.Pub(&foo{"bar!"})

	time.Sleep(5 * time.Second)
}
