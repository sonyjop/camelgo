package core

type Component interface {
	GetScheme() string
	CreateEndpoint(EndpointConfig) (Endpoint, error)
}
type EndpointConfig struct {
	RawURI    string
	Scheme    string
	Params    map[string]interface{}
	Component Component
}
type Endpoint interface {
	CreateProducer() (Producer, error)
	CreateConsumer(target Processor) (Consumer, error)
	GetURI() string
}

type Producer interface {
	Processor
	Start(ctx Context) error
	Stop(ctx Context) error
}

type Consumer interface {
	Start(ctx Context) error
	Stop(ctx Context) error
}
