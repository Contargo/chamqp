package exchange_declare

import (
	"github.com/Contargo/chamqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompleteCustomDecl(t *testing.T) {
	t.Run("fills decl correctly", func(t *testing.T) {
		expectedSpec := chamqp.ExchangeDeclareSpec{
			Name:       "testme",
			Kind:       "fanout",
			Durable:    true,
			AutoDelete: false,
			Internal:   true,
			NoWait:     false,
			Args:       amqp.Table{},
			ErrorChan:  nil,
		}
		spec := DeclareExchange("testme").
			WithKind("fanout").
			WithDurable(true).
			WithAutoDelete(false).
			WithInternal(true).
			WithNoWait(false).
			WithArgs(amqp.Table{}).
			WithErrorChan(nil).
			BuildSpec()
		assert.Equal(t, expectedSpec, spec)
	})
}

func TestWithCertainDefaults(t *testing.T) {
	t.Run("custom Exchange, default Kind", func(t *testing.T) {
		DeclareExchange("testme").
			WithDefaultKind().
			Defaults().
			BuildSpec()
	})

	t.Run("custom Exchange, Kind, default Durable", func(t *testing.T) {
		DeclareExchange("testme").
			WithKind("fan").
			WithDefaultDurable().
			Defaults().
			BuildSpec()
	})

	t.Run("custom Exchange, Kind, Durable, default autodelete", func(t *testing.T) {
		DeclareExchange("testme").
			WithKind("fan").
			WithDurable(false).
			WithDefaultAutoDelete().
			Defaults().
			BuildSpec()
	})

	t.Run("custom Exchange, Kind, Durable, autodelete, default Internal", func(t *testing.T) {
		DeclareExchange("testme").
			WithKind("fan").
			WithDurable(false).
			WithAutoDelete(false).
			WithDefaultInternal().
			Defaults().
			BuildSpec()
	})
	t.Run("custom Exchange, Kind, Durable, autodelete, Internal", func(t *testing.T) {
		DeclareExchange("testme").
			WithKind("fan").
			WithDurable(false).
			WithAutoDelete(false).
			WithInternal(false).
			WithDefaultNoWait().
			Defaults().
			BuildSpec()
	})
}

func TestWithAllDefaults(t *testing.T) {
	t.Run("with Name custom", func(t *testing.T) {
		DeclareExchange("").
			Defaults().
			BuildSpec()
	})

	t.Run("with Name custom, custom Kind", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "custom"
		expectedSpec.Kind = "asdf"

		r := DeclareExchange("custom").
			WithKind("asdf").
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("with Name custom, custom Kind, custom Durable", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "asdf"
		expectedSpec.Kind = "asdf"
		expectedSpec.Durable = true

		r := DeclareExchange("asdf").
			WithKind("asdf").
			WithDurable(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custom Name, Kind, Durable, AutoDelete", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "testme"
		expectedSpec.Kind = "asdf"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = false

		r := DeclareExchange("testme").
			WithKind("asdf").
			WithDurable(true).
			WithAutoDelete(false).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})

	t.Run("custon Name, Kind, Durable, autodelete, Internal", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "testme"
		expectedSpec.Kind = "asdf"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = false
		expectedSpec.Internal = true

		r := DeclareExchange("testme").
			WithKind("asdf").
			WithDurable(true).
			WithAutoDelete(false).
			WithInternal(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
	t.Run("custon Name, Kind, Durable, autodelete, Internal, nowait", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "testme"
		expectedSpec.Kind = "asdf"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = false
		expectedSpec.Internal = true
		expectedSpec.NoWait = true

		r := DeclareExchange("testme").
			WithKind("asdf").
			WithDurable(true).
			WithAutoDelete(false).
			WithInternal(true).
			WithNoWait(true).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
	t.Run("custon Name, Kind, Durable, autodelete, Internal, nowait, Args", func(t *testing.T) {
		expectedSpec := Defaults
		expectedSpec.Name = "testme"
		expectedSpec.Kind = "asdf"
		expectedSpec.Durable = true
		expectedSpec.AutoDelete = false
		expectedSpec.Internal = true
		expectedSpec.NoWait = true
		expectedSpec.Args = amqp.Table{}

		r := DeclareExchange("testme").
			WithKind("asdf").
			WithDurable(true).
			WithAutoDelete(false).
			WithInternal(true).
			WithNoWait(true).
			WithArgs(amqp.Table{}).
			Defaults().
			BuildSpec()
		assert.Equal(t, expectedSpec, r)
	})
}
