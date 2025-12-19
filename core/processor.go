package core

import "context"

// Processor is the universal runtime contract.
type Processor interface {
	Process(ctx context.Context, exchange *Exchange) error
}

// PipelineProcessor runs steps sequentially.
type PipelineProcessor struct {
	Children []Processor
}

func (p *PipelineProcessor) Process(ctx context.Context, exchange *Exchange) error {
	for _, child := range p.Children {
		if exchange.Error != nil {
			return exchange.Error
		}
		if err := child.Process(ctx, exchange); err != nil {
			return err
		}
		// Handle In/Out promotion here
	}
	return nil
}

// ChoiceProcessor handles Content-Based Routing.
type ChoiceProcessor struct {
	Branches []struct {
		Condition Predicate
		Pipeline  Processor
	}
	Otherwise Processor
}
