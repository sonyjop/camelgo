package core

// RouteLoader is the bridge between a source and the engine's IR.
type RouteLoader interface {
	// Load takes a source (e.g., a RouteBuilder instance or a file path)
	// and returns the Intermediate Representation (RouteDefinitions).
	Load(source interface{}) ([]*RouteDefinition, error)
}
