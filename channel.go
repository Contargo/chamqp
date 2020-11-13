package chamqp

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type consumeSpec struct {
	queue     string
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp.Table

	deliveryChan chan<- amqp.Delivery
	errorChan    chan<- error
}

type exchangeDeclareSpec struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp.Table

	errorChan chan<- error
}

type queueBindSpec struct {
	name     string
	key      string
	exchange string
	noWait   bool
	args     amqp.Table

	errorChan chan<- error
}

type queueDeclareSpec struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table

	queueChan chan<- amqp.Queue
	errorChan chan<- error
}

// Channel represents an AMQP channel. Used as a context for valid message
// exchange. Errors on methods with this Channel will be detected and the
// channel will recreate itself.
type Channel struct {
	ch                   *amqp.Channel
	consumeSpecs         []consumeSpec
	exchangeDeclareSpecs []exchangeDeclareSpec
	queueBindSpecs       []queueBindSpec
	queueDeclareSpecs    []queueDeclareSpec
	mu                   sync.Mutex
}

func (ch *Channel) connected(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		ch.ch = nil
		return err
	}
	ch.ch = channel

	for _, spec := range ch.exchangeDeclareSpecs {
		ch.applyExchangeDeclareSpec(spec)
	}
	for _, spec := range ch.queueDeclareSpecs {
		ch.applyQueueDeclareSpec(spec)
	}
	for _, spec := range ch.queueBindSpecs {
		ch.applyQueueBindSpec(spec)
	}
	for _, spec := range ch.consumeSpecs {
		ch.applyConsumeSpec(spec)
	}

	return nil
}

func (ch *Channel) disconnected() {
	ch.ch = nil
}

func (ch *Channel) applyExchangeDeclareSpec(spec exchangeDeclareSpec) {
	err := ch.ch.ExchangeDeclare(spec.name, spec.kind, spec.durable, spec.autoDelete, spec.internal, spec.noWait, spec.args)
	if err != nil && spec.errorChan != nil {
		spec.errorChan <- err
	}
}

func (ch *Channel) applyQueueDeclareSpec(spec queueDeclareSpec) {
	queue, err := ch.ch.QueueDeclare(spec.name, spec.durable, spec.autoDelete, spec.exclusive, spec.noWait, spec.args)
	if err != nil && spec.errorChan != nil {
		spec.errorChan <- err
		return
	}
	if spec.queueChan != nil {
		spec.queueChan <- queue
	}
}

func (ch *Channel) applyQueueBindSpec(spec queueBindSpec) {
	err := ch.ch.QueueBind(spec.name, spec.key, spec.exchange, spec.noWait, spec.args)
	if err != nil && spec.errorChan != nil {
		spec.errorChan <- err
	}
}

func (ch *Channel) applyConsumeSpec(spec consumeSpec) {
	deliveries, err := ch.ch.Consume(spec.queue, spec.consumer, spec.autoAck, spec.exclusive, spec.noLocal, spec.noWait, spec.args)
	if err != nil && spec.errorChan != nil {
		spec.errorChan <- err
		return
	}
	if spec.deliveryChan != nil {
		go shovel(deliveries, spec.deliveryChan)
	}
}

// Consume immediately starts delivering queued messages.
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table, deliveryChan chan<- amqp.Delivery, errorChan chan<- error) {
	spec := consumeSpec{
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
		deliveryChan,
		errorChan,
	}
	ch.consumeSpecs = append(ch.consumeSpecs, spec)
	if ch.ch != nil {
		ch.applyConsumeSpec(spec)
	}
}

// Publish sends a Publishing from the client to an exchange on the server.
func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if ch.ch == nil {
		return fmt.Errorf("context has no channel")
	}

	return ch.ch.Publish(exchange, key, mandatory, immediate, msg)
}

// ExchangeDeclare declares an exchange on the server. If the exchange does not
// already exist, the server will create it. If the exchange exists, the server
// verifies that it is of the provided type, durability and auto-delete flags.
func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table, errorChan chan<- error) {
	spec := exchangeDeclareSpec{
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
		errorChan,
	}
	ch.exchangeDeclareSpecs = append(ch.exchangeDeclareSpecs, spec)
	if ch.ch != nil {
		ch.applyExchangeDeclareSpec(spec)
	}
}

// QueueBind binds an exchange to a queue so that publishings to the exchange
// will be routed to the queue when the publishing routing key matches the
// binding routing key.
func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table, errorChan chan<- error) {
	spec := queueBindSpec{
		name,
		key,
		exchange,
		noWait,
		args,
		errorChan,
	}
	ch.queueBindSpecs = append(ch.queueBindSpecs, spec)
	if ch.ch != nil {
		ch.applyQueueBindSpec(spec)
	}
}

// QueueDeclare declares a queue to hold messages and deliver to consumers.
// Declaring creates a queue if it doesn't already exist, or ensures that an
// existing queue matches the same parameters.
func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table, queueChan chan<- amqp.Queue, errorChan chan<- error) {
	spec := queueDeclareSpec{
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
		queueChan,
		errorChan,
	}
	ch.queueDeclareSpecs = append(ch.queueDeclareSpecs, spec)
	if ch.ch != nil {
		ch.applyQueueDeclareSpec(spec)
	}
}

// Shovel takes messages from `src` and puts them into `dest`.
func shovel(src <-chan amqp.Delivery, dest chan<- amqp.Delivery) {
	for msg := range src {
		dest <- msg
	}
}
