package core

type Message struct {
	body    interface{}
	headers map[string]interface{}
}

func NewMessage() *Message {
	return &Message{
		headers: make(map[string]interface{}),
	}
}
func (m *Message) Body() interface{} {
	return m.body
}
func (m *Message) SetBody(body interface{}) {
	m.body = body
}
func (m *Message) Headers() map[string]interface{} {
	return m.headers
}
func (m *Message) SetHeader(key string, value interface{}) {
	m.headers[key] = value
}
func (m *Message) Header(key string) interface{} {
	return m.headers[key]
}

type Exchange struct {
	id         string
	in         *Message
	out        *Message
	err        error
	properties map[string]interface{}
}

func NewExchange() *Exchange {
	return &Exchange{
		id:         "1",
		in:         NewMessage(),
		out:        NewMessage(),
		properties: make(map[string]interface{}),
	}
}
func (e *Exchange) ID() string {
	return e.id
}
func (e *Exchange) In() *Message {
	return e.in
}
func (e *Exchange) SetIn(msg *Message) {
	e.in = msg
}
func (e *Exchange) Out() *Message {
	return e.out
}
func (e *Exchange) SetOut(msg *Message) {
	e.out = msg
}
func (e *Exchange) Properties() map[string]interface{} {
	return e.properties
}
func (e *Exchange) SetProperty(key string, value interface{}) {
	e.properties[key] = value
}
func (e *Exchange) GetProperty(key string) interface{} {
	return e.properties[key]
}
func (e *Exchange) SetError(err error) {
	e.err = err
}
func (e *Exchange) Error() error {
	return e.err
}
func (e *Exchange) Clone() Exchange {
	clone := &Exchange{
		id:         e.id + "_clone",
		in:         &Message{body: e.In().Body(), headers: mapCloner(e.In().Headers())},
		out:        &Message{body: e.Out().Body(), headers: mapCloner(e.Out().Headers())},
		properties: make(map[string]interface{}),
	}
	for k, v := range e.properties {
		clone.properties[k] = v
	}
	return *clone
}
func mapCloner(original map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for k, v := range original {
		clone[k] = v
	}
	return clone
}

type Predicate interface {
	Evaluate(ctx Context, exchange *Exchange) (bool, error)
}
