package consume

import (
	"github.com/Contargo/chamqp"
	"github.com/streadway/amqp"
)

var Defaults = chamqp.ConsumeSpec{
	Queue:        "",
	DeliveryChan: nil,

	Consumer:  "",
	AutoAck:   true,
	Exclusive: false,
	NoLocal:   false,
	NoWait:    false,
	Args:      nil,
	ErrorChan: nil,
}

func Consume(queueName string) QueueDecl {
	consumeDecl := Defaults
	consumeDecl.Queue = queueName
	return QueueDecl{consumeDecl}
}

type QueueDecl struct {
	consumeSpec chamqp.ConsumeSpec
}

func (q QueueDecl) WithDeliveryChan(deliveryChan chan amqp.Delivery) DeliveryChan {
	q.consumeSpec.DeliveryChan = deliveryChan
	return DeliveryChan{q}
}

type DeliveryChan struct {
	queueDecl QueueDecl
}

func (d DeliveryChan) Defaults() ErrorChan {
	return ErrorChan{d.queueDecl}
}

func (d DeliveryChan) WithConsumer(consumer string) ConsumerDecl {
	d.queueDecl.consumeSpec.Consumer = consumer
	return ConsumerDecl{d.queueDecl}
}

func (d DeliveryChan) WithDefaultConsumer() ConsumerDecl {
	return ConsumerDecl{d.queueDecl}
}

type ConsumerDecl struct {
	queueDecl QueueDecl
}

func (d ConsumerDecl) WithAutoAck(autoAck bool) AutoAckDecl {
	d.queueDecl.consumeSpec.AutoAck = autoAck
	return AutoAckDecl{d.queueDecl}
}

func (d ConsumerDecl) WithDefaultAutoAck() AutoAckDecl {
	return AutoAckDecl{d.queueDecl}
}

func (d ConsumerDecl) Defaults() ErrorChan {
	return ErrorChan{d.queueDecl}
}

type AutoAckDecl struct {
	queueDecl QueueDecl
}

func (a AutoAckDecl) Defaults() ErrorChan {
	return ErrorChan{a.queueDecl}
}

func (a AutoAckDecl) WithExclusive(exclusive bool) ExclusiveDecl {
	a.queueDecl.consumeSpec.Exclusive = false
	return ExclusiveDecl{a.queueDecl}
}

func (a AutoAckDecl) WithDefaultExclusive() ExclusiveDecl {
	return ExclusiveDecl{a.queueDecl}
}

type ExclusiveDecl struct {
	queueDecl QueueDecl
}

func (e ExclusiveDecl) Defaults() ErrorChan {
	return ErrorChan{e.queueDecl}
}

func (e ExclusiveDecl) WithNoLocal(noLocal bool) NoLocalDecl {
	e.queueDecl.consumeSpec.NoLocal = noLocal
	return NoLocalDecl{e.queueDecl}
}

func (e ExclusiveDecl) WithDefaultNoLocal() NoWaitDecl {
	return NoLocalDecl{e.queueDecl}
}

type NoLocalDecl struct {
	queueDecl QueueDecl
}

func (n NoLocalDecl) Defaults() ErrorChan {
	return ErrorChan{n.queueDecl}
}

func (n NoLocalDecl) WithNoWait(noWait bool) NoWaitDecl {
	n.queueDecl.consumeSpec.NoWait = noWait
	return NoWaitDecl{n.queueDecl}
}

type NoWaitDecl struct {
	queueDecl QueueDecl
}

func (n NoWaitDecl) Defaults() ErrorChan {
	return ErrorChan{n.queueDecl}
}

func (n NoWaitDecl) WithArgs(args amqp.Table) ArgsDecl {
	n.queueDecl.consumeSpec.Args = args
	return ArgsDecl{n.queueDecl}
}

type ArgsDecl struct {
	queueDecl QueueDecl
}

func (a ArgsDecl) Defaults() ErrorChan {
	return ErrorChan{a.queueDecl}
}

func (a ArgsDecl) WithErrorChan(errorChan chan error) ErrorChan {
	a.queueDecl.consumeSpec.ErrorChan = errorChan
	return ErrorChan{a.queueDecl}
}

type ErrorChan struct {
	queueDecl QueueDecl
}

func (e ErrorChan) BuildSpec() chamqp.ConsumeSpec {
	return e.queueDecl.consumeSpec
}

func (e ErrorChan) Build(ch *chamqp.Channel) {
	ch.ConsumeWithSpec(e.queueDecl.consumeSpec)
}
