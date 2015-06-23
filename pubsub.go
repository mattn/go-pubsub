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
	"runtime"
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

func NewWrap(f interface{}) *wrap {
	return &wrap{f: f}
}

// PubSub contains channel and callbacks.
type PubSub struct {
	c chan interface{}
	f []interface{}
	w []*wrap
	m sync.Mutex
	e chan error
}

// New return new PubSub intreface.
func New() *PubSub {
	ps := new(PubSub)
	ps.c = make(chan interface{})
	ps.e = make(chan error)
	call := func(f interface{}, rf reflect.Value, in []reflect.Value) {
		defer func() {
			if err := recover(); err != nil {
				ps.e <- &PubSubError{f, err}
			}
		}()
		rf.Call(in)
	}

	go func() {
		for v := range ps.c {
			rv := reflect.ValueOf(v)
			ps.m.Lock()
			for _, f := range ps.f {
				rf := reflect.ValueOf(f)
				if rv.Type() == reflect.ValueOf(f).Type().In(0) {
					go call(f, rf, []reflect.Value{rv})
				}
			}
			for _, w := range ps.w { // wrapped functions
				rf := reflect.ValueOf(w.f)
				if rv.Type() == reflect.ValueOf(w.f).Type().In(0) {
					go call(w.f, rf, []reflect.Value{rv})
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
func (ps *PubSub) Sub(f interface{}) error {
	check := f
	w, wrapped := f.(*wrap)
	if wrapped { // check wrapped function instead
		check = w.f
	}
	rf := reflect.ValueOf(check)
	if rf.Kind() != reflect.Func {
		return errors.New("Not a function")
	}
	if rf.Type().NumIn() != 1 {
		return errors.New("Number of arguments should be 1")
	}
	ps.m.Lock()
	defer ps.m.Unlock()
	if wrapped { // append to wrapped set
		ps.w = append(ps.w, w)
	} else {
		ps.f = append(ps.f, f)
	}
	return nil
}

// Leave unsubscribe to the PubSub.
func (ps *PubSub) Leave(f interface{}) {
	var fp uintptr
	if f == nil {
		if pc, _, _, ok := runtime.Caller(1); ok {
			fp = runtime.FuncForPC(pc).Entry()
		}
	} else {
		fp = reflect.ValueOf(f).Pointer()
	}
	ps.m.Lock()
	defer ps.m.Unlock()
	if _, wrapped := f.(*wrap); wrapped {
		result := make([]*wrap, 0, len(ps.w))
		last := 0
		for i, v := range ps.w {
			if reflect.ValueOf(v).Pointer() == fp {
				result = append(result, ps.w[last:i]...)
				last = i + 1
			}
		}
		ps.w = append(result, ps.w[last:]...)
	} else {
		result := make([]interface{}, 0, len(ps.f))
		last := 0
		for i, v := range ps.f {
			if reflect.ValueOf(v).Pointer() == fp {
				result = append(result, ps.f[last:i]...)
				last = i + 1
			}
		}
		ps.f = append(result, ps.f[last:]...)
	}
}

// Pub publish to the PubSub.
func (ps *PubSub) Pub(v interface{}) {
	ps.c <- v
}

// Close closes PubSub. To inspect unbsubscribing for another subscruber, you must create message structure to notify them. After publish notifycations, Close should be called.
func (ps *PubSub) Close() {
	close(ps.c)
	ps.f = nil
	ps.w = nil
}
