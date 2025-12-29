package core

// CompileContext is the minimal context needed to compile definitions.
// It decouples compilation from the full Context API.
type CompileContext interface {
	GetEndpoint(uri string) (Endpoint, error)
	NewExchange() *Exchange
}

// Compilable is a metadata node that knows how to turn itself into a Processor.
type Compilable interface {
	Compile(ctx CompileContext) (Processor, error)
}

// RouteDefinition is the top-level blueprint.
type RouteDefinition struct {
	ID       string
	InputURI string
	Steps    []Compilable // The IR tree
}
