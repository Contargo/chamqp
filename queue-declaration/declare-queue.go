package queue_declaration

import (
	"github.com/streadway/amqp"
	"gitlab.contargo.net/karrei/chamqp"
	"gitlab.contargo.net/karrei/chamqp/queue-bind"
)

var defaults = chamqp.QueueDeclareSpec{
	Name: "",
	Durable: false,
	AutoDelete: true,
	Exclusive: false,
	NoWait: false,
	Args: amqp.Table{},
	QueueChan: nil,
	ErrorChan: nil,
}

func DeclareQueue(name string) NameDecl {
	return DeclareQueueWithChan(name, nil)
}

func DeclareQueueWithChan(name string, amqpChan *amqp.Channel) NameDecl {
	queueDecl := defaults
	queueDecl.Name = name
	return NameDecl{queueDecl, amqpChan}
}

type NameDecl struct {
	queueDecl chamqp.QueueDeclareSpec
	amqpChan *amqp.Channel
}

func (n NameDecl) WithDurable(durable bool) DurableDecl {
	n.queueDecl.Durable = durable
	return DurableDecl{n}
}

func (n NameDecl) Defaults() End {
	return End{n}
}

type DurableDecl struct {
	nameDecl NameDecl
}

func (d DurableDecl) WithAutoDelete(autodelete bool) AutoDelete {
	d.nameDecl.queueDecl.AutoDelete = autodelete
	return AutoDelete{d.nameDecl}
}

func (d DurableDecl) Defaults() End {
	return End{d.nameDecl}
}

type AutoDelete struct {
	nameDecl NameDecl
}

func (a AutoDelete) WithExclusive(exclusive bool) ExclusiveDecl {
	a.nameDecl.queueDecl.Exclusive = exclusive
	return ExclusiveDecl{a.nameDecl}
}

func (a AutoDelete) Defaults() End {
	return End{a.nameDecl}
}

type ExclusiveDecl struct {
	nameDecl NameDecl
}

func (e ExclusiveDecl) WithNoWait(noWait bool) NoWaitDecl {
	e.nameDecl.queueDecl.NoWait = noWait
	return NoWaitDecl{e.nameDecl}
}

func (e ExclusiveDecl) Defaults() End {
	return End{e.nameDecl}
}

type NoWaitDecl struct {
	nameDecl NameDecl
}

func (n NoWaitDecl) WithArgs(args amqp.Table) ArgsDecl {
	n.nameDecl.queueDecl.Args = args
	return ArgsDecl{n.nameDecl}
}

func (n NoWaitDecl) Defaults() End {
	return End{n.nameDecl}
}

type ArgsDecl struct {
	nameDecl NameDecl
}

func (a ArgsDecl) WithQueueChan(queueChan chan amqp.Queue) QueueChanDecl {
	a.nameDecl.queueDecl.QueueChan = queueChan
	return QueueChanDecl{a.nameDecl}
}

func (a ArgsDecl) Defaults() End {
	return End{a.nameDecl}
}

type QueueChanDecl struct {
	nameDecl NameDecl
}

func (q QueueChanDecl) Defaults() End {
	return End{q.nameDecl}
}

func (q QueueChanDecl) WithErrorChan(errorChan chan error) End {
	q.nameDecl.queueDecl.ErrorChan = errorChan
	return End{q.nameDecl}
}

type End struct {
	nameDecl NameDecl
}

func (e End) BuildSpec() chamqp.QueueDeclareSpec{
	return e.nameDecl.queueDecl
}

func (e End) Build(ch *chamqp.Channel) {
	ch.QueueDeclareWithSpec(e.nameDecl.queueDecl)
}

func (e End) AndBind() queue_bind.NameDecl {
	return queue_bind.BindQueue(e.nameDecl.queueDecl.Name)
}