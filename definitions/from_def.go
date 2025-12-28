package definitions

import "github.com/sonyjop/camelgo/core"

// FromDefinition is the metadata for consuming from an endpoint.
// It represents the input source of a route step (not the route's InputURI).
// When compiled, it creates a Consumer that feeds messages into the next processor.
type FromDefinition struct {
	URI string
}

func (d *FromDefinition) Compile(ctx core.CompileContext) (core.Processor, error) {
	ep, err := ctx.GetEndpoint(d.URI)
	if err != nil {
		return nil, err
	}
	// Create a consumer that wraps the next processor
	// Note: the actual pipeline processor will be passed separately during route compilation
	// For now, return a placeholder processor that would be enhanced with a consumer
	cons, err := ep.CreateConsumer(nil)
	if err != nil {
		return nil, err
	}
	// Consumer implements the start/stop lifecycle; wrap it as a Processor if needed
	// This is a simplified version; in a real scenario you might need a ConsumerProcessor adapter
	return &FromProcessor{Consumer: cons}, nil
}

// FromProcessor adapts a Consumer to the Processor interface
type FromProcessor struct {
	Consumer core.Consumer
}

func (fp *FromProcessor) Process(ctx core.Context, exchange *core.Exchange) error {
	// FromProcessor is primarily a lifecycle holder; actual message flow is handled by Consumer
	// This Process method is a no-op as the Consumer pushes messages independently
	return nil
}
