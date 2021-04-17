package chamqp

import "github.com/Contargo/chamqp/publish"

func PublishWithChannel(channel *Channel) publish.WithExchange {
	return &publish.PublishStruct{
		false,
		false,
		"",
		"",
		channel,
	}
}

