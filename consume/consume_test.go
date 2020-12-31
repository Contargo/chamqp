package consume

import "testing"

func TestCustomDeclare(t *testing.T) {
	t.Run("fill decls complete", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			WithExclusive(true).
			WithNoLocal(true).
			WithNoWait(true).
			WithArgs(nil).
			WithErrorChan(nil).
			BuildSpec()
	})
}

func TestWithDefaults(t *testing.T) {
	t.Run("custom queue, deliveryChan", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			Defaults().
			BuildSpec()
	})
	
	t.Run("custom queue, deliveryChan, consumer", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			Defaults().
			BuildSpec()
	})
	
	t.Run("custom queue, deliveryChan, consumer, autoack", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			Defaults().
			BuildSpec()
	})

	t.Run("custom queue, deliveryChan, consumer, autoack, exclusive", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			WithExclusive(true).
			Defaults().
			BuildSpec()
	})

	t.Run("custom queue, deliveryChan, consumer, autoack, exclusive, nolocal", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			WithExclusive(true).
			WithNoLocal(true).
			Defaults().
			BuildSpec()
	})

	t.Run("custom queue, deliveryChan, consumer, autoack, exclusive, nolocal, nowait", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			WithExclusive(true).
			WithNoLocal(true).
			WithNoWait(true).
			Defaults().
			BuildSpec()
	})

	t.Run("custom queue, deliveryChan, consumer, autoack, exclusive, nolocal, nowait, args", func(t *testing.T) {
		Consume("testqueue").
			WithDeliveryChan(nil).
			WithConsumer("consumer").
			WithAutoAck(false).
			WithExclusive(true).
			WithNoLocal(true).
			WithNoWait(true).
			WithArgs(nil).
			Defaults().
			BuildSpec()
	})
}
