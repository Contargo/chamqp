package queue_bind

import (
	"github.com/Contargo/chamqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomDeclare(t *testing.T) {
	t.Run("fills decls correctly", func(t *testing.T) {
		expectedSpec := chamqp.QueueBindSpec{
			Name:      "testqueue",
			Key:       "routingKey",
			Exchange:  "exchange",
			NoWait:    false,
			Args:      nil,
			ErrorChan: nil,
		}

		r := BindQueue("testqueue").
			WithExchange("exchange").
			WithRoutingKey("routingKey").
			WithNoWait(false).
			WithArgs(nil).
			WithErrorChannel(nil).
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
}

func TestWithDefaults(t *testing.T) {
	t.Run("custom name, routingkey, exchange", func(t *testing.T) {
		expectedSpec := chamqp.QueueBindSpec{
			Name:      "test",
			Key:       "routing",
			Exchange:  "exchangeName",
			NoWait:    false,
			Args:      nil,
			ErrorChan: nil,
		}
		r := BindQueue("test").
			WithExchange("exchangeName").
			WithRoutingKey("routing").
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
}
