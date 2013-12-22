package multichannel

import (
	"errors"
	"reflect"
	"runtime"
)

type MultiChannel struct {
	c chan interface{}
	f []interface{}
}

func New() *MultiChannel {
	mc := new(MultiChannel)
	mc.c = make(chan interface{})
	go func() {
		for v := range mc.c {
			rv := reflect.ValueOf(v)
			for _, f := range mc.f {
				rf := reflect.ValueOf(f)
				if rv.Type() == reflect.ValueOf(f).Type().In(0) {
					rf.Call([]reflect.Value{rv})
				}
			}
		}
	}()
	return mc
}

func (mc *MultiChannel) Sub(f interface{}) error {
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func {
		return errors.New("Not a function")
	}
	if rf.Type().NumIn() != 1 {
		return errors.New("Number of arguments should be 1")
	}
	mc.f = append(mc.f, f)
	return nil
}

func (mc *MultiChannel) Leave(f interface{}) {
	var fp uintptr
	if f == nil {
		if pc, _, _, ok := runtime.Caller(1); ok {
			fp = runtime.FuncForPC(pc).Entry()
		}
	} else {
		fp = reflect.ValueOf(f).Pointer()
	}
    result := make([]interface{}, 0, len(mc.f))
    last := 0
    for i, v := range mc.f {
        if reflect.ValueOf(v).Pointer() == fp {
            result = append(result, mc.f[last:i]...)
            last = i + 1
        }
    }
    mc.f = append(result, mc.f[last:]...)
}

func (mc *MultiChannel) Pub(v interface{}) {
	mc.c <-v
}
