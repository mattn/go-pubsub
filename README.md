go-pubsub
=========

PubSub model message hub

Usage
-----

To subscribe:

    ps := pubsub.New()
    ps.Sub(func(i int) {
        fmt.Println("int subscriber: ", i)
    })

To publish:

    ps.Pub(1)

Messages are passed to subscriber which have same type argument.

License
-------

MIT: http://mattn.mit-license.org/2013

Author
------

Yasuhiro Matsumoto (mattn.jp@gmail.com)
