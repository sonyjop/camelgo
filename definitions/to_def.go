package definitions

import "github.com/sonyjop/camelgo/core"

// ToDefinition is the metadata for sending to an endpoint.
type ToDefinition struct {
	URI string
}

func (d *ToDefinition) Compile(ctx core.CompileContext) (core.Processor, error) {
	ep, _ := ctx.GetEndpoint(d.URI)
	prod, _ := ep.CreateProducer()
	return prod, nil // Producer implements Processor
}
