package queue_declaration

import (
	"github.com/Contargo/chamqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomDeclare(t *testing.T) {
	t.Run("fills decls correctly", func(t *testing.T) {
		expectedSpec := chamqp.QueueDeclareSpec{
			Name:       "queue",
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       amqp.Table{},
			QueueChan:  nil,
			ErrorChan:  nil,
		}
		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(false).
			WithExclusive(false).
			WithNoWait(false).
			WithArgs(amqp.Table{}).
			WithQueueChan(nil).
			WithErrorChan(nil).
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
}

func TestWithDefaults(t *testing.T) {
	t.Run("custom name", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"

		r := DeclareQueue("queue").
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom name, durable", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true

		r := DeclareQueue("queue").
			WithDurable(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom name, durable, autodelete", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = true

		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom queue, durable, autodelete, exclusive", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = true
		expectedSpec.Exclusive = true

		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(true).
			WithExclusive(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom queue, durable, autodelete, exclusive", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = true
		expectedSpec.Exclusive = true
		expectedSpec.NoWait = true

		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(true).
			WithExclusive(true).
			WithNoWait(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom queue, durable, autodelete, exclusive", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = true
		expectedSpec.Exclusive = true
		expectedSpec.NoWait = true
		expectedSpec.Args = nil

		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(true).
			WithExclusive(true).
			WithNoWait(true).
			WithArgs(nil).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom queue, durable, autodelete, exclusive", func(t *testing.T) {
		expectedSpec := defaults
		expectedSpec.Name = "queue"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = true
		expectedSpec.Exclusive = true
		expectedSpec.NoWait = true
		expectedSpec.Args = nil
		expectedSpec.QueueChan = nil

		r := DeclareQueue("queue").
			WithDurable(true).
			WithAutoDelete(true).
			WithExclusive(true).
			WithNoWait(true).
			WithArgs(nil).
			WithQueueChan(nil).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

}
