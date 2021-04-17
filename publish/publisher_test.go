package publish

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ChannelMock struct {
	immediate      bool
	mandatory      bool
	key            string
	exchangeName   string
	objectToBeSent interface{}
}

func (c *ChannelMock) PublishJSON(exchange, key string, mandatory, immediate bool, objectToBeSent interface{}) error {
	c.immediate = immediate
	c.mandatory = mandatory
	c.key = key
	c.exchangeName = exchange
	c.objectToBeSent = objectToBeSent
	return nil
}

func NewChannelMock() ChannelWithPublishJson{
	return &ChannelMock{}
}

func TestPublishing(t *testing.T) {
	t.Run("fill custom data correctly", func(t *testing.T) {
		channel := NewChannelMock()
		callData := &ChannelMock{
			true,
			true,
			"asdf",
			"123",
			"aaa",
		}
		
		WithChannel(channel).
			WithExchange(callData.exchangeName).
			WithRoutingKey(callData.key).
			WithMandatory(callData.mandatory).
			WithImmediate(callData.immediate).
			Publish(callData.objectToBeSent)
		
		assert.Equal(t,callData,channel)
	})

	t.Run("fill defaultdata correctly", func(t *testing.T) {
		channel := NewChannelMock()
		callData := &ChannelMock{
			false,
			false,
			"bbb",
			"ccc",
			"ddd",
		}

		WithChannel(channel).
			WithExchange(callData.exchangeName).
			WithRoutingKey(callData.key).
			Publish(callData.objectToBeSent)

		assert.Equal(t,callData,channel)
	})
}
