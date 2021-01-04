package queue_bind

import (
	"github.com/stretchr/testify/assert"
	"gitlab.contargo.net/karrei/chamqp"
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
			WithRoutinghKey("routingKey").
			WithExchangeDecl("exchange").
			WithNoWaitDecl(false).
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
			NoWait:    true,
			Args:      nil,
			ErrorChan: nil,
		}
		r := BindQueue("test").
			WithRoutinghKey("routing").
			WithExchangeDecl("exchangeName").
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
}
