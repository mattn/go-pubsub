package pubsub_test

import (
	"github.com/mattn/go-pubsub"
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	done := make(chan int)
	ps := pubsub.New()
	ps.Sub(func(i int) {
		done <- i
	})
	ps.Pub(1)
	i := <-done
	if i != 1 {
		t.Fatalf("Expected %v, but %d:", 1, i)
	}
}

func TestString(t *testing.T) {
	done := make(chan string)
	ps := pubsub.New()
	ps.Sub(func(s string) {
		done <- s
	})
	ps.Pub("hello world")
	s := <-done
	if s != "hello world" {
		t.Fatalf("Expected %v, but %d:", "hello world", s)
	}
}

type F struct {
	m string
}

func TestStruct(t *testing.T) {
	done := make(chan *F)
	ps := pubsub.New()
	ps.Sub(func(f *F) {
		done <- f
	})
	ps.Pub(&F{"hello world"})
	f := <-done
	if f.m != "hello world" {
		t.Fatalf("Expected %v, but %d:", "hello world", f.m)
	}
}

func TestOnly(t *testing.T) {
	doneInt := make(chan int)
	doneF := make(chan *F)
	ps := pubsub.New()
	ps.Sub(func(i int) {
		doneInt <- i
	})
	ps.Sub(func(f *F) {
		doneF <- f
	})
	ps.Pub(&F{"hello world"})
	ps.Pub(2)
	i := <-doneInt
	f := <-doneF
	if f.m != "hello world" {
		t.Fatalf("Expected %v, but %d:", "hello world", f.m)
	}
	if i != 2 {
		t.Fatalf("Expected %v, but %d:", 2, f.m)
	}
}

func TestClojure(t *testing.T) {
	done := make(chan int)
	ps := pubsub.New()
	add := func() {
		ps.Sub(func(i int) {
			done <- i
		})
	}
	add()
	add()
	ps.Pub(1)

	time.AfterFunc(3*time.Second, func() {
		close(done)
	})
	i1 := <-done
	i2 := <-done
	if i1 != 1 || i2 != 1 {
		t.Fatal("Expected multiple subscribers")
	}
}

func TestLeave(t *testing.T) {
	done := make(chan int)
	ps := pubsub.New()

	f := func(i int) {
		done <- i
	}
	ps.Sub(f)
	ps.Sub(f)
	ps.Pub(1)
	i1 := <-done
	i2 := <-done
	if i1 != 1 || i2 != 1 {
		t.Fatal("Expected multiple subscribers")
	}
	ps.Leave(f)
	ps.Pub(2)
	select {
	case <-done:
		t.Fatal("WTF")
	default:
	}
}

func TestWrap(t *testing.T) {
	done := make(chan int)
	ps := pubsub.New()

	f := func(i int) {
		done <- i
	}
	w1, w2 := pubsub.NewWrap(f), pubsub.NewWrap(f)
	ps.Sub(w1)
	ps.Sub(w2)
	ps.Pub(1)
	i1 := <-done
	i2 := <-done
	if i1 != 1 || i2 != 1 {
		t.Fatal("Expected multiple subscribers")
	}
	ps.Leave(w2)
	ps.Pub(2)
	i1 = <-done
	select {
	case i2 = <-done:
		t.Fatal("WTF")
	default:
	}
	if i1 != 2 {
		t.Fatal("Expected single subscribers")
	}
	ps.Leave(w1)
	ps.Pub(3)
	select {
	case <-done:
		t.Fatal("WTF")
	default:
	}
}
