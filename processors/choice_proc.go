package processors

import "github.com/sonyjop/camelgo/core"

// ChoiceProcessor handles Content-Based Routing.
type ChoiceProcessor struct {
	Branches []struct {
		Condition core.Processor
	}
	Otherwise core.Processor
}

func (c *ChoiceProcessor) Process(ctx core.Context, exchange *core.Exchange) error {
	/*for _, branch := range c.Branches {
		match, err := branch.Condition.Evaluate(exchange)
		if err != nil {
			return err
		}
		if match {
			return branch.Pipeline.Process(ctx, exchange)
		}
	}
	if c.Otherwise != nil {
		return c.Otherwise.Process(ctx, exchange)
	}*/
	return nil
}
