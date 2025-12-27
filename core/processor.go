package core

// Processor is the universal runtime contract.
type Processor interface {
	Process(ctx Context, exchange *Exchange) error
}

// PipelineProcessor runs steps sequentially.
type PipelineProcessor struct {
	Children []Processor
}

func (p *PipelineProcessor) Process(ctx Context, exchange *Exchange) error {
	for _, child := range p.Children {
		if exchange.Error() != nil {
			return exchange.Error()
		}
		if err := child.Process(ctx, exchange); err != nil {
			return err
		}
		// Handle In/Out promotion here
	}
	return nil
}
