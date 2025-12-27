package core

// Definition is a metadata node that knows how to turn itself into a Processor.
type Definition interface {
	Compile(ctx Context) (Processor, error)
}

// RouteDefinition is the top-level blueprint.
type RouteDefinition struct {
	ID       string
	InputURI string
	Steps    []Definition // The IR tree
}
