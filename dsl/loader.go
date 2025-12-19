package dsl

import (
	"fmt"

	"github.com/sonyjop/camelgo/core"
)

// DSLLoader is the implementation of core.RouteLoader for Go DSL.
type DSLLoader struct {
	// We might need a reference to context to help initialize definitions
}

func NewDSLLoader() *DSLLoader {
	return &DSLLoader{}
}

// Load processes a RouteBuilder to generate RouteDefinitions.
func (l *DSLLoader) Load(source interface{}) ([]*core.RouteDefinition, error) {
	// 1. Type check: Ensure the source is a RouteBuilder
	builder, ok := source.(core.RouteBuilder)
	if !ok {
		return nil, fmt.Errorf("DSLLoader expected core.RouteBuilder, got %T", source)
	}

	// 2. Execute the user's DSL configuration
	// This populates the internal state of the builder
	builder.Configure()

	// 3. Extract the resulting IR (Blueprints)
	// We assume RouteBuilder has a GetRouteDefinitions() method
	definitions := builder.GetRouteDefinitions()

	if len(definitions) == 0 {
		return nil, fmt.Errorf("no routes were defined in the provided RouteBuilder")
	}

	return definitions, nil
}
