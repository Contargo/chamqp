package publish_json_wait

import "time"

type PublishJsonAndWaitForResponseStruct struct {
	replyQueueName  string
	correlationId   string
	response        interface{}
	request         interface{}
	exchange        string
	key             string
	mandatory       bool
	immediate       bool
	responseTimeout time.Duration
	t               PublishJsonAndWaitForResponseContract
}
type ReplyQueueName interface {
	WithReplyQueueName(replyQueueName string) CorrelationId
}
type CorrelationId interface {
	WithCorrelationId(correlationId string) Response
}
type Response interface {
	WithResponse(response interface{}) Request
}
type Request interface {
	WithRequest(request interface{}) Exchange
}
type Exchange interface {
	WithExchange(exchange string) Key
}
type Key interface {
	WithKey(key string) Mandatory
}
type Mandatory interface {
	PublishJsonAndWaitForResponse
	WithMandatory(mandatory bool) Immediate
}
type Immediate interface {
	WithImmediate(immediate bool) ResponseTimeout
}
type ResponseTimeout interface {
	WithResponseTimeout(responseTimeout time.Duration) PublishJsonAndWaitForResponse
}
type PublishJsonAndWaitForResponse interface {
	PublishJsonAndWaitForResponse() error
}

func (t *PublishJsonAndWaitForResponseStruct) WithReplyQueueName(replyQueueName string) CorrelationId {
	t.replyQueueName = replyQueueName
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithCorrelationId(correlationId string) Response {
	t.correlationId = correlationId
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithResponse(response interface{}) Request {
	t.response = response
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithRequest(request interface{}) Exchange {
	t.request = request
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithExchange(exchange string) Key {
	t.exchange = exchange
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithKey(key string) Mandatory {
	t.key = key
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithMandatory(mandatory bool) Immediate {
	t.mandatory = mandatory
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithImmediate(immediate bool) ResponseTimeout {
	t.immediate = immediate
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) WithResponseTimeout(responseTimeout time.Duration) PublishJsonAndWaitForResponse {
	t.responseTimeout = responseTimeout
	return t
}
func (t *PublishJsonAndWaitForResponseStruct) PublishJsonAndWaitForResponse() error {
	return t.t.PublishJsonAndWaitForResponse(t.replyQueueName, t.correlationId, t.response, t.request, t.exchange, t.key, t.mandatory, t.immediate, t.responseTimeout)
}
func WithChannel(t PublishJsonAndWaitForResponseContract) ReplyQueueName {
	ret := PublishJsonAndWaitForResponseStruct{}
	ret.responseTimeout = 10 * time.Second
	ret.mandatory = false
	ret.immediate = false
	ret.t = t
	return &ret
}

type PublishJsonAndWaitForResponseContract interface {
	PublishJsonAndWaitForResponse(replyQueueName string, correlationId string, response interface{}, request interface{}, exchange string, key string, mandatory bool, immediate bool, responseTimeout time.Duration) error
}
