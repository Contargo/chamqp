package chamqp

import (
	"math"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	initialInterval = 1 * time.Second
	maxInterval     = 10 * time.Second
	multiplier      = float64(2)
)

// Connection manages the serialization and deserialization of frames from IO
// and dispatches the frames to the appropriate channel. All RPC methods and
// asynchronous Publishing, Delivery, Ack, Nack and Return messages are
// multiplexed on this channel. There must always be active receivers for every
// asynchronous message on this connection.
type Connection struct {
	conn                   *amqp.Connection
	channels               []*Channel
	errorChans             []chan error
	shutdownChan, doneChan chan struct{}
	mu                     sync.Mutex
}

// Dial accepts a string in the AMQP URI format and returns a new Connection
// over TCP using PlainAuth. Defaults to a server heartbeat interval of 10
// seconds and sets the handshake deadline to 30 seconds. After handshake,
// deadlines are cleared.
//
// Use `NotifyError` to register a receiver for errors on the connection.
func Dial(url string) *Connection {
	conn := &Connection{
		shutdownChan: make(chan struct{}),
		doneChan:     make(chan struct{}),
	}
	go conn.supervise(url)
	return conn
}

func (c *Connection) connect(url string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	c.conn = conn

	for _, ctx := range c.channels {
		ctx.connected(c.conn)
	}

	return nil
}

func (c *Connection) disconnect(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		for _, c := range c.errorChans {
			c <- err
		}
	}

	c.conn = nil

	for _, ctx := range c.channels {
		ctx.disconnected()
	}
}

func (c *Connection) supervise(url string) {
	var attempt float64

	defer close(c.doneChan)

	for {
		backoffDelay := time.Duration(math.Pow(multiplier, attempt)) * initialInterval
		if backoffDelay > maxInterval {
			backoffDelay = maxInterval
		}

		err := c.connect(url)
		if err != nil {
			for _, c := range c.errorChans {
				c <- err
			}
			attempt++
			select {
			case <-time.After(backoffDelay):
				continue
			case <-c.shutdownChan:
				return
			}
		}
		attempt = 0

		notifyClose := make(chan *amqp.Error)
		c.conn.NotifyClose(notifyClose)

		select {
		case err := <-notifyClose:
			c.disconnect(err)
		case <-c.shutdownChan:
			return
		}
	}
}

// NotifyError registers a listener for error events either initiated by an
// connect or close.
func (c *Connection) NotifyError(receiver chan error) chan error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.errorChans = append(c.errorChans, receiver)

	return receiver
}

// Channel opens a unique, concurrent server channel to process the bulk of AMQP
// messages. Any error from methods on this receiver will cause the Channel to
// recreate itself.
func (c *Connection) Channel() *Channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := &Channel{}

	c.channels = append(c.channels, ch)

	if c.conn != nil {
		ch.connected(c.conn)
	}

	return ch
}

// Close requests and waits for the response to close the AMQP connection.
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.shutdownChan:
		// Already closed. Nothing to do.
	default:
		close(c.shutdownChan)
	}

	<-c.doneChan

	if c.conn != nil {
		conn := c.conn
		c.conn = nil
		return conn.Close()
	}
	return nil
}
