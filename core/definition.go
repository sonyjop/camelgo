package core

//import "github.com/sonyjop/camelgo/core"

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

// ToDefinition is the metadata for sending to an endpoint.
type ToDefinition struct {
	URI string
}

func (d *ToDefinition) Compile(ctx Context) (Processor, error) {
	ep, _ := ctx.GetEndpoint(d.URI)
	prod, _ := ep.CreateProducer()
	return prod, nil // Producer implements Processor
}

// ChoiceDefinition is the metadata for a branching flow.
type ChoiceDefinition struct {
	WhenClauses []WhenDefinition
	Otherwise   []Definition
}

func (d *ChoiceDefinition) Compile(ctx Context) (Processor, error) {
	// Logic to recursively call Compile() on all branches and return a ChoiceProcessor
}
