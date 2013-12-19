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
		fmt.Println("foo subscriber: ", f.bar)
	})

	mc.Pub(1)
	mc.Pub("hello")
	mc.Pub(2)
	mc.Pub(&foo{"bar!"})

	time.Sleep(5 * time.Second)
}
