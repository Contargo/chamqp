package queue_bind

import (
	"github.com/Contargo/chamqp"
	"github.com/Contargo/chamqp/consume"
	amqp "github.com/rabbitmq/amqp091-go"
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

func BindQueueWithExchange(queueName, exchange string) ExchangeDecl {
	bind := Defaults
	bind.Name = queueName
	bind.Exchange = exchange
	return ExchangeDecl{NameDecl{bind}}
}

type NameDecl struct {
	queueBindSpec chamqp.QueueBindSpec
}

func (n ExchangeDecl) WithRoutingKey(routingKey string) KeyDecl {
	n.nameDecl.queueBindSpec.Key = routingKey
	return KeyDecl{n.nameDecl}
}

type KeyDecl struct {
	nameDecl NameDecl
}

func (k NameDecl) WithExchange(exchangeName string) ExchangeDecl {
	k.queueBindSpec.Exchange = exchangeName
	return ExchangeDecl{k}
}

type ExchangeDecl struct {
	nameDecl NameDecl
}

func (e KeyDecl) Defaults() ErrorChanDecl {
	return ErrorChanDecl{e.nameDecl}
}

func (e KeyDecl) WithNoWait(noWait bool) NoWaitDecl {
	e.nameDecl.queueBindSpec.NoWait = noWait
	return NoWaitDecl{e.nameDecl}
}

func (e KeyDecl) WithDefaultNoWait() NoWaitDecl {
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

func (n NoWaitDecl) WithDefaultArgs() ArgsDecl {
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

func (a ArgsDecl) WithDefaultErrorChannel() ErrorChanDecl {
	return ErrorChanDecl{a.nameDecl}
}

type ErrorChanDecl struct {
	nameDecl NameDecl
}

func (e ErrorChanDecl) BuildSpec() chamqp.QueueBindSpec {
	return e.nameDecl.queueBindSpec
}

type ConsumeDecl struct {
	nameDecl NameDecl
}

func (e ErrorChanDecl) Build(channel *chamqp.Channel) ConsumeDecl {
	channel.QueueBindWithSpec(e.nameDecl.queueBindSpec)
	return ConsumeDecl{e.nameDecl}
}

func (c ConsumeDecl) AndConsume() consume.QueueDecl {
	return consume.Consume(c.nameDecl.queueBindSpec.Name)
}
