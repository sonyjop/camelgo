package core

//import "context"

// Route is the "Live" version of a RouteDefinition.
// It is created during the 'Reification' (Compilation) phase.
type Route struct {
	ID       string
	InputURI string

	// Consumer is the active listener (e.g., KafkaConsumer, HttpServer)
	Consumer Consumer

	// Pipeline is the top-level Processor that contains all the logic
	Pipeline Processor

	// Reference back to the context for resource access
	context Context
}

// Start activates the consumer to begin receiving messages.
func (r *Route) Start(ctx Context) error {
	if r.Consumer == nil {
		return nil // Or return an error if a route must have a consumer
	}
	return r.Consumer.Start(ctx)
}

// Stop gracefully shuts down the consumer.
func (r *Route) Stop(ctx Context) error {
	if r.Consumer == nil {
		return nil
	}
	return r.Consumer.Stop(ctx)
}
