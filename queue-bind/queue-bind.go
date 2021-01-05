package queue_bind

import (
	"github.com/Contargo/chamqp"
	"github.com/Contargo/chamqp/consume"
	"github.com/streadway/amqp"
)

var Defaults = chamqp.QueueBindSpec{
	Name:      "",
	Key:       "",
	Exchange:  "",
	NoWait:    false,
	Args:      nil,
	ErrorChan: nil,
}

func BindQueue(queueName string) NameDecl {
	bind := Defaults
	bind.Name = queueName
	return NameDecl{bind}
}

type NameDecl struct {
	queueBindSpec chamqp.QueueBindSpec
}

func (n NameDecl) WithRoutinghKey(routingKey string) KeyDecl {
	n.queueBindSpec.Key = routingKey
	return KeyDecl{n}
}

type KeyDecl struct {
	nameDecl NameDecl
}

func (k KeyDecl) WithExchangeDecl(exchangeName string) ExchangeDecl {
	k.nameDecl.queueBindSpec.Exchange = exchangeName
	return ExchangeDecl{k.nameDecl}
}

type ExchangeDecl struct {
	nameDecl NameDecl
}

func (e ExchangeDecl) Defaults() ErrorChanDecl {
	return ErrorChanDecl{e.nameDecl}
}

func (e ExchangeDecl) WithNoWaitDecl(noWait bool) NoWaitDecl {
	e.nameDecl.queueBindSpec.NoWait = noWait
	return NoWaitDecl{e.nameDecl}
}

type NoWaitDecl struct {
	nameDecl NameDecl
}

func (n NoWaitDecl) Defaults() ErrorChanDecl {
	return ErrorChanDecl{n.nameDecl}
}

func (n NoWaitDecl) WithArgs(args amqp.Table) ArgsDecl {
	n.nameDecl.queueBindSpec.Args = args
	return ArgsDecl{n.nameDecl}
}

type ArgsDecl struct {
	nameDecl NameDecl
}

func (a ArgsDecl) Defaults() ErrorChanDecl {
	return ErrorChanDecl{a.nameDecl}
}

func (a ArgsDecl) WithErrorChannel(channel chan error) ErrorChanDecl {
	a.nameDecl.queueBindSpec.ErrorChan = channel
	return ErrorChanDecl{a.nameDecl}
}

type ErrorChanDecl struct {
	nameDecl NameDecl
}

func (e ErrorChanDecl) BuildSpec() chamqp.QueueBindSpec {
	return e.nameDecl.queueBindSpec
}

func (e ErrorChanDecl) Build(channel *chamqp.Channel) {
	channel.QueueBindWithSpec(e.nameDecl.queueBindSpec)
}

func (e ErrorChanDecl) AndConsume() consume.QueueDecl {
	return consume.Consume(e.nameDecl.queueBindSpec.Name)
}
