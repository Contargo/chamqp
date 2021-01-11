package combined

import (
	"github.com/Contargo/chamqp"
	"github.com/Contargo/chamqp/exchange-declare"
	"github.com/Contargo/chamqp/queue-declaration"
	"testing"
)

func TestCombined(t *testing.T) {
	t.Run("Test combined declare, bind and consume", func(t *testing.T) {
		channel := &chamqp.Channel{}
		queue_declaration.
			DeclareQueue("MySuperQueue").	Defaults().
			Build(channel).
			AndBind().
			WithExchange("myExchange").	
			WithRoutingKey("myroutingKey").Defaults().
			Build(channel).
			AndConsume().
			WithDeliveryChan(nil).Defaults().
			Build(channel)
	})
	t.Run("Test combined exchange declare, queue declare, bind and consume", func(t *testing.T) {
		channel := &chamqp.Channel{}
		exchange_declare.DeclareExchange("myExchange").
			Defaults().
			Build(channel).
			AndDeclareQueue("MySuperQueue").
			Defaults().
			Build(channel).
			AndBindWithExchange().
			WithRoutingKey("#").
			Defaults().
			Build(channel).
			AndConsume().
			WithDeliveryChan(nil).
			Defaults().
			Build(channel)
		
	})	
}
