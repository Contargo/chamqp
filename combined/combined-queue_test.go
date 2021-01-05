package combined

import (
	"github.com/Contargo/chamqp/queue-declaration"
	"testing"
)

func TestCombined(t *testing.T) {
	t.Run("Test combined declare, bind and consume", func(t *testing.T) {
		queue_declaration.
			DeclareQueue("MySuperQueue").	Defaults().
			AndBind().
			WithRoutinghKey("myroutingKey").
			WithExchange("myExchange").	Defaults().
			AndConsume().
			WithDeliveryChan(nil).Defaults().
			BuildSpec()
	})
}
