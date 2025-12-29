package core

// RouteBuilder is the interface for DSL builders that can configure routes.
type RouteBuilder interface {
	Configure()
	GetRouteDefinitions() []*RouteDefinition
}

// RouteLoader is the bridge between a source and the engine's IR.
type RouteLoader interface {
	// Load takes a source (e.g., a RouteBuilder instance or a file path)
	// and returns the Intermediate Representation (RouteDefinitions).
	Load(source interface{}) ([]*RouteDefinition, error)
}
