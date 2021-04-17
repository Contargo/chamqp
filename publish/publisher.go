package publish

type ChannelWithPublishJson interface {
	PublishJSON(exchange, key string, mandatory, immediate bool, objectToBeSent interface{}) error
}

type Publish interface {
	Publish(objectToBeSent interface{})
}

type Mandatory interface {
	WithMandatory(mandatory bool) Immediate
	Publish
}

type Immediate interface {
	WithImmediate(bool) Publish
	Publish
}

type WithRoutingKey interface {
	WithRoutingKey(key string) Mandatory
}

type WithExchange interface {
	WithExchange(exchangeName string) WithRoutingKey
}

type PublishStruct struct {
	immediate    bool
	mandatory    bool
	key          string
	exchangeName string
	channel      ChannelWithPublishJson
}

func (p *PublishStruct) Publish(objectToBeSent interface{}) {
	p.channel.PublishJSON(p.exchangeName, p.key, p.mandatory, p.immediate, objectToBeSent)
}

func (p *PublishStruct) WithImmediate(immediate bool) Publish {
	p.immediate = immediate
	return p
}

func (p *PublishStruct) WithMandatory(mandatory bool) Immediate {
	p.mandatory = mandatory
	return p
}

func (p *PublishStruct) WithRoutingKey(key string) Mandatory {
	p.key = key
	return p
}

func (p *PublishStruct) WithExchange(exchangeName string) WithRoutingKey {
	p.exchangeName = exchangeName
	return p
}

func WithChannel(channel ChannelWithPublishJson) WithExchange {
	return &PublishStruct{
		false,
		false,
		"",
		"",
		channel,
	}
}
