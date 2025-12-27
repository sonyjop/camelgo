package core

type Component interface {
	GetScheme() string
	CreateEndpoint(uri string, options map[string]interface{}) (Endpoint, error)
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
