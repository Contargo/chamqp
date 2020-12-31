package combined

import (
	"gitlab.contargo.net/karrei/chamqp/queue-declaration"
	"testing"
)

func TestCombined(t *testing.T) {
	t.Run("Test combined declare, bind and consume", func(t *testing.T) {
		queue_declaration.
			DeclareQueue("MySuperQueue").	Defaults().
			AndBind().
			WitRoutinghKey("myroutingKey").
			WithExchangeDecl("myExchange").	Defaults().
			AndConsume().
			WithDeliveryChan(nil).Defaults().
			BuildSpec()
	})
}
