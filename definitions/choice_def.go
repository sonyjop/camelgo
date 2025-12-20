package definitions

import (
	"github.com/sonyjop/camelgo/core"
	"github.com/sonyjop/camelgo/processors"
)

// ChoiceDefinition holds the blueprint for branching logic
type ChoiceDefinition struct {
	WhenClauses []WhenDefinition
	Otherwise   []core.Definition
}

type WhenDefinition struct {
	Condition core.Predicate
	Steps     []core.Definition
}

// Compile transforms the IR into a ChoiceProcessor
func (d *ChoiceDefinition) Compile(ctx core.Context) (core.Processor, error) {
	runtimeChoice := &processors.ChoiceProcessor{}

	return runtimeChoice, nil
}
