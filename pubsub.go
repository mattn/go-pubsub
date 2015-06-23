/*
This package provides a portable intaface to pubsub model.
PubSub can publish/subscribe/unsubscribe messages for all.
To subscribe:

	ps := pubsub.New()
	ps.Sub(func(s string) {
		fmt.Println(s)
	})

To publish:

	ps.Pub("hello world")

The message are allowed to pass any types, and passing to subscribers which
can accept the type for the argument of callback.
*/
package pubsub

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type PubSubError struct {
	f interface{}
	e interface{}
}

func (pse *PubSubError) String() string {
	return fmt.Sprintf("%v: %v", pse.f, pse.e)
}

func (pse *PubSubError) Error() string {
	return fmt.Sprint(pse.e)
}

func (pse *PubSubError) Subscriber() interface{} {
	return pse.f
}

type wrap struct {
	f interface{}
}

func NewFunc(f interface{}) *wrap {
	w := &wrap{}
	w.f = f
	return w
}

// PubSub contains channel and callbacks.
type PubSub struct {
	c chan interface{}
	w []*wrap
	m sync.Mutex
	e chan error
}

// New return new PubSub intreface.
func New() *PubSub {
	ps := new(PubSub)
	ps.c = make(chan interface{})
	ps.e = make(chan error)
	go func() {
		for v := range ps.c {
			rv := reflect.ValueOf(v)
			ps.m.Lock()
			for _, w := range ps.w {
				rf := reflect.ValueOf(w.f)
				if rv.Type() == reflect.ValueOf(w.f).Type().In(0) {
					go func(f interface{}, rf reflect.Value) {
						defer func() {
							if err := recover(); err != nil {
								ps.e <- &PubSubError{f, err}
							}
						}()
						rf.Call([]reflect.Value{rv})
					}(w.f, rf)
				}
			}
			ps.m.Unlock()
		}
	}()
	return ps
}

func (ps *PubSub) Error() chan error {
	return ps.e
}

// Sub subscribe to the PubSub.
func (ps *PubSub) Sub(w *wrap) error {
	f := w.f
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func {
		return errors.New("Not a function")
	}
	if rf.Type().NumIn() != 1 {
		return errors.New("Number of arguments should be 1")
	}
	ps.m.Lock()
	defer ps.m.Unlock()
	ps.w = append(ps.w, w)
	return nil
}

// Leave unsubscribe to the PubSub.
func (ps *PubSub) Leave(w *wrap) {
	ps.m.Lock()
	defer ps.m.Unlock()
	for k, v := range ps.w {
		if w == v {
			ps.w = append(ps.w[:k], ps.w[k+1:]...)
			return
		}
	}
}

// Pub publish to the PubSub.
func (ps *PubSub) Pub(v interface{}) {
	ps.c <- v
}

// Close closes PubSub. To inspect unbsubscribing for another subscruber, you must create message structure to notify them. After publish notifycations, Close should be called.
func (ps *PubSub) Close() {
	close(ps.c)
	ps.w = nil
}
