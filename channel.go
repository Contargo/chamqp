package chamqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotifyPublishSpec struct {
	confirm chan amqp.Confirmation
}

type ConsumeSpec struct {
	Queue        string
	Consumer     string
	DeliveryChan chan<- amqp.Delivery

	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
	ErrorChan chan<- error
}

type ExchangeDeclareSpec struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table

	ErrorChan chan<- error
}

type QueueBindSpec struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table

	ErrorChan chan<- error
}

type QueueDeclareSpec struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table

	QueueChan chan<- amqp.Queue
	ErrorChan chan<- error
}

type Properties struct {
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id
}

// Channel represents an AMQP channel. Used as a context for valid message
// Exchange. Errors on methods with this Channel will be detected and the
// channel will recreate itself.
type Channel struct {
	ch                   *amqp.Channel
	consumeSpecs         []ConsumeSpec
	exchangeDeclareSpecs []ExchangeDeclareSpec
	queueBindSpecs       []QueueBindSpec
	queueDeclareSpecs    []QueueDeclareSpec
	notifyPublishSpec    []NotifyPublishSpec
	mu                   sync.Mutex
	confirm              bool
	confirmNoWait        bool
}

func (ch *Channel) connected(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if ch.confirm {
		channel.Confirm(ch.confirmNoWait)
	}
	if err != nil {
		ch.ch = nil
		return err
	}
	ch.ch = channel

	for _, spec := range ch.exchangeDeclareSpecs {
		err := ch.applyExchangeDeclareSpec(spec)
		if err != nil {
			return err
		}
	}
	for _, spec := range ch.queueDeclareSpecs {
		err := ch.applyQueueDeclareSpec(spec)
		if err != nil {
			return err
		}
	}
	for _, spec := range ch.queueBindSpecs {
		err := ch.applyQueueBindSpec(spec)
		if err != nil {
			return err
		}
	}
	for _, spec := range ch.consumeSpecs {
		err := ch.applyConsumeSpec(spec)
		if err != nil {
			return err
		}
	}
	for _, spec := range ch.notifyPublishSpec {
		ch.applyNotifyPublishSpec(spec)
	}

	return nil
}

func (ch *Channel) disconnected() {
	ch.ch = nil
}

func (ch *Channel) applyExchangeDeclareSpec(spec ExchangeDeclareSpec) error {
	err := ch.ch.ExchangeDeclare(spec.Name, spec.Kind, spec.Durable, spec.AutoDelete, spec.Internal, spec.NoWait, spec.Args)
	if err != nil {
		if spec.ErrorChan != nil {
			spec.ErrorChan <- err
		}
		return err
	}
	return nil
}

func (ch *Channel) applyQueueDeclareSpec(spec QueueDeclareSpec) error {
	queue, err := ch.ch.QueueDeclare(spec.Name, spec.Durable, spec.AutoDelete, spec.Exclusive, spec.NoWait, spec.Args)
	if err != nil {
		if spec.ErrorChan != nil {
			spec.ErrorChan <- err
		}
		return err
	}
	if spec.QueueChan != nil {
		spec.QueueChan <- queue
	}
	return nil
}

func (ch *Channel) applyQueueBindSpec(spec QueueBindSpec) error {
	err := ch.ch.QueueBind(spec.Name, spec.Key, spec.Exchange, spec.NoWait, spec.Args)
	if err != nil {
		if spec.ErrorChan != nil {
			spec.ErrorChan <- err
		}
		spec.ErrorChan <- err
		return err
	}
	return nil
}

func (ch *Channel) applyConsumeSpec(spec ConsumeSpec) error {
	deliveries, err := ch.ch.Consume(spec.Queue, spec.Consumer, spec.AutoAck, spec.Exclusive, spec.NoLocal, spec.NoWait, spec.Args)
	if err != nil {
		if spec.ErrorChan != nil {
			spec.ErrorChan <- err
		}
		return err
	}
	if spec.DeliveryChan != nil {
		go shovel(deliveries, spec.DeliveryChan)
	}
	return nil
}

func (ch *Channel) applyNotifyPublishSpec(spec NotifyPublishSpec) {
	subscribeChannel := make(chan amqp.Confirmation, 1)
	go shovelConfirmation(subscribeChannel, spec.confirm)
	ch.ch.NotifyPublish(subscribeChannel)
}

func (ch *Channel) ConsumeWithSpec(spec ConsumeSpec) {
	ch.Consume(spec.Queue, spec.Consumer, spec.AutoAck, spec.Exclusive, spec.NoLocal, spec.NoWait, spec.Args, spec.DeliveryChan, spec.ErrorChan)
}

// Consume immediately starts delivering queued messages.
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table, deliveryChan chan<- amqp.Delivery, errorChan chan<- error) {
	spec := ConsumeSpec{
		queue,
		consumer,
		deliveryChan,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
		errorChan,
	}
	ch.consumeSpecs = append(ch.consumeSpecs, spec)
	if ch.ch != nil {
		ch.applyConsumeSpec(spec)
	}
}

// Publish sends a Publishing from the client to an Exchange on the server.
func (ch *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if ch.ch == nil {
		return fmt.Errorf("context has no channel")
	}

	return ch.ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (ch *Channel) PublishJSONWithProperties(exchange, key string, mandatory, immediate bool, objectToBeSent interface{}, properties Properties) error {
	if ch.ch == nil {
		return fmt.Errorf("context has no channel")
	}

	payload, err := json.Marshal(objectToBeSent)
	if err != nil {
		return err
	}
	return ch.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,

		ContentEncoding: properties.ContentEncoding,
		DeliveryMode:    properties.DeliveryMode,
		Priority:        properties.Priority,
		CorrelationId:   properties.CorrelationId,
		ReplyTo:         properties.ReplyTo,
		Expiration:      properties.Expiration,
		MessageId:       properties.MessageId,
		Timestamp:       properties.Timestamp,
		Type:            properties.Type,
		UserId:          properties.UserId,
		AppId:           properties.AppId,
	})
}

func (ch *Channel) PublishJSON(exchange, key string, mandatory, immediate bool, objectToBeSent interface{}) error {
	if ch.ch == nil {
		return fmt.Errorf("context has no channel")
	}

	payload, err := json.Marshal(objectToBeSent)
	if err != nil {
		return err
	}
	return ch.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
}

func (ch *Channel) PublishJsonAndWaitForResponse(replyQueueName, correlationId string, response, request interface{}, exchange, key string, mandatory, immediate bool, responseTimeout time.Duration) error {
	if ch.ch == nil {
		return errors.New("channel not present")
	}
	defer ch.ch.Cancel(replyQueueName+".consumer", false)
	replyQueue, err := ch.ch.Consume(replyQueueName, replyQueueName+".consumer", true, false, false, false, nil)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers:       amqp.Table{},
		ContentType:   "application/json",
		ReplyTo:       replyQueueName,
		CorrelationId: correlationId,
		Body:          payload,
	}
	err = ch.Publish(exchange, key, mandatory, immediate, msg)
	if err != nil {
		return err
	}

	timer := time.NewTimer(responseTimeout)
	for {
		select {
		case reply := <-replyQueue:
			if correlationId != "" && reply.CorrelationId != correlationId {
				fmt.Println("skipping, not for me")
				continue
			}
			err := json.Unmarshal(reply.Body, response)
			if err != nil {
				continue
			}
			return err
		case <-timer.C:
			return errors.New("timed out")
		}
	}
}

func (ch *Channel) ExchangeDeclareWithSpec(spec ExchangeDeclareSpec) {
	ch.ExchangeDeclare(spec.Name, spec.Kind, spec.Durable, spec.AutoDelete, spec.Internal, spec.NoWait, spec.Args, spec.ErrorChan)
}

// ExchangeDeclare declares an Exchange on the server. If the Exchange does not
// already exist, the server will create it. If the Exchange exists, the server
// verifies that it is of the provided type, durability and auto-delete flags.
func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table, errorChan chan<- error) {
	spec := ExchangeDeclareSpec{
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

func (ch *Channel) QueueBindWithSpec(q QueueBindSpec) {
	ch.QueueBind(q.Name, q.Key, q.Exchange, q.NoWait, q.Args, q.ErrorChan)
}

// QueueBind binds an Exchange to a Queue so that publishings to the Exchange
// will be routed to the Queue when the publishing routing Key matches the
// binding routing Key.
func (ch *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table, errorChan chan<- error) {
	spec := QueueBindSpec{
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

func (ch *Channel) QueueDeclareWithSpec(q QueueDeclareSpec) {
	ch.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Args, q.QueueChan, q.ErrorChan)
}

// QueueDeclare declares a Queue to hold messages and deliver to consumers.
// Declaring creates a Queue if it doesn't already exist, or ensures that an
// existing Queue matches the same parameters.
func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table, queueChan chan<- amqp.Queue, errorChan chan<- error) {
	spec := QueueDeclareSpec{
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

func (ch *Channel) NotifyPublish() chan amqp.Confirmation {
	notifyPublishChan := make(chan amqp.Confirmation, 1)
	spec := NotifyPublishSpec{notifyPublishChan}
	ch.notifyPublishSpec = append(ch.notifyPublishSpec, spec)
	if ch.ch != nil {
		ch.applyNotifyPublishSpec(spec)
	}
	return notifyPublishChan
}

// Shovel takes messages from `src` and puts them into `dest`.
func shovel(src <-chan amqp.Delivery, dest chan<- amqp.Delivery) {
	for msg := range src {
		dest <- msg
	}
}

func shovelConfirmation(src, dest chan amqp.Confirmation) {
	for msg := range src {
		dest <- msg
	}
}
