package main

import (
	"fmt"
	"github.com/mattn/go-multichannel"
	"time"
)

func main() {
	mc := multichannel.New()
	mc.Sub(func(i int) {
		fmt.Println("int subscriber: ", i)
	})
	mc.Sub(func(s string) {
		fmt.Println("string subscriber: ", s)
	})

	mc.Pub(1)
	mc.Pub("hello")
	mc.Pub(2)

	time.Sleep(5 * time.Second)
}
