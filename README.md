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

# Getting started for development

Simply clone this repository
```sh
$ git clone https://github.com/Contargo/chamqp
```


