Chamqp
======


# Features
Chamqp is a small layer above github.com/streadway/amqp featuring auto-reconnect and limited exponential back-off. This is especially useful when running multiple services and network disconnect will happen sooner or later in production.

## Examples

# Plans for the future
The API for creating exchanges is overloaded with (still important parameters). A builder pattern will be provided later for easier handling. Also we plan to add at least some integration tests.

# Getting started

Simply run in your project
```sh
$ go get github.com/Contargo/chamqp 
```

# Usage - classic way
Chamqp is built with the intention to be compatible to the underlying github.com/streadway/amqp package. 
This means instead of amqp.Dial() you have to use chamqp.Dial(). On the connection itself you then can use .Channel(). 
Declaring queues or exchanges is done in the same fashion as with streadway's amqp package.


See the following example to ilustrate it:
```
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

# Usage with builder
It's cumbersome to keep track about all the paramters especially when multiple boolean flags are used. Therefor we added a small implementation using the builder pattern.

```
BindQueue("testqueue").
    WithRoutinghKey("routingKey").
    WithExchangeDecl("exchange").
    WithNoWaitDecl(false).
    WithArgs(nil).
    WithErrorChannel(nil).
    Build(channel)
```

You can also use default values so you don't have to type everything:

```
BindQueue("test").
    WithRoutinghKey("routing").
    WithExchangeDecl("exchangeName").
    Defaults().
    BuildSpec()
```

For further samples have a look at the _test.go files

Defaults are held in a public accessible variable:
* queue_bind.Defaults
* exchange-declare.Defaults
* consume.Defaults

# Getting started for development

Simply clone this repository
```sh
$ git clone https://github.com/Contargo/chamqp
```


