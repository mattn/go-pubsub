go-pubsub
=========

PubSub model message hub

Usage
-----

To subscribe:

```go
ps := pubsub.New()
ps.Sub(func(i int) {
    fmt.Println("int subscriber: ", i)
})
```

To publish:

```go
ps.Pub(1)
```

If the closure captures values, you can use a wrapped function, eg:

```go
// global
var ps = pubsub.New()
// subscribe with a network connection
func Subscribe(conn net.Conn) {
	f := func(i int) { // a closure captures "conn", it will send messages to different network connections
		conn.Write([]byte(fmt.Sprint("int subscriber: ", i)))
	})
	ps.Sub(pubsub.NewWrap(f)) // subscribe with a wrapper
	...
}
```

Messages are passed to subscriber which have same type argument.

License
-------

MIT: http://mattn.mit-license.org/2013

Author
------

Yasuhiro Matsumoto (mattn.jp@gmail.com)
