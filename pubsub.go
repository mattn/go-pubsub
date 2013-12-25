package pubsub

import (
	"errors"
	"reflect"
	"runtime"
)

type PubSub struct {
	c chan interface{}
	f []interface{}
}

func New() *PubSub {
	ps := new(PubSub)
	ps.c = make(chan interface{})
	go func() {
		for v := range ps.c {
			rv := reflect.ValueOf(v)
			for _, f := range ps.f {
				rf := reflect.ValueOf(f)
				if rv.Type() == reflect.ValueOf(f).Type().In(0) {
					rf.Call([]reflect.Value{rv})
				}
			}
		}
	}()
	return ps
}

func (ps *PubSub) Sub(f interface{}) error {
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func {
		return errors.New("Not a function")
	}
	if rf.Type().NumIn() != 1 {
		return errors.New("Number of arguments should be 1")
	}
	ps.f = append(ps.f, f)
	return nil
}

func (ps *PubSub) Leave(f interface{}) {
	var fp uintptr
	if f == nil {
		if pc, _, _, ok := runtime.Caller(1); ok {
			fp = runtime.FuncForPC(pc).Entry()
		}
	} else {
		fp = reflect.ValueOf(f).Pointer()
	}
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

func (ps *PubSub) Pub(v interface{}) {
	ps.c <-v
}
