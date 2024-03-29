# Chamqp


## Features

Chamqp is a small layer above [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) featuring auto-reconnect using exponential back-off with an upper bound. This is especially useful when running multiple services and network disconnect will happen sooner or later in production.


## Getting started

Simply run in your project
```sh
$ go get github.com/Contargo/chamqp
```


## Usage - classic way

Chamqp is built with the intention to be compatible to the underlying [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) package.
This means instead of `amqp.Dial()` you have to use `chamqp.Dial()`. On the connection itself you then can use `.Channel()`.
Declaring queues or exchanges is done in the same fashion as with the `amqp091-go` package.

See the following example to illustrate it:

```go
conn, err := chamqp.Dial(applicatonConfig.AMQPUrl)
channel := conn.Channel()
channel.ExchangeDeclare("exchangeName", "topic", false, false, false, false, nil, errChan)
...

channel.Publish(
    "exchangeName,
    "routing.key",
    false,
    false,
    amqp.Publishing{
        ContentType: "contentType",
        Body: "payload",
    },
)
```


## Usage with builder

Experimental - use at your own risk.

It's cumbersome to keep track about all the parameters especially when multiple boolean flags are used. Therefor we added a small implementation using the builder pattern.

```go
BindQueue("testqueue").
    WithRoutinghKey("routingKey").
    WithExchangeDecl("exchange").
    WithNoWaitDecl(false).
    WithArgs(nil).
    WithErrorChannel(nil).
    Build(channel)
```

You can also use default values, so you don't have to type everything:

```go
BindQueue("test").
    WithRoutinghKey("routing").
    WithExchangeDecl("exchangeName").
    Defaults().
    BuildSpec()
```

For further samples have a look at the `_test.go` files

Defaults are held in a public accessible variable:
* queue_bind.Defaults
* exchange-declare.Defaults
* consume.Defaults


## Getting started for development

Simply clone this repository
```sh
$ git clone https://github.com/Contargo/chamqp
```


