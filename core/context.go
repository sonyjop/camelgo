package core

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// Context is the interface that components and processors interact with.
type Context interface {
	// Lifecycle
	Start() error
	Stop() error

	// Route Management
	SetLoader(loader RouteLoader)
	AddRoutes(source interface{}) error

	// Registry & Factory
	RegisterComponent(scheme string, component Component)
	GetComponent(scheme string) (Component, error)
	GetEndpoint(uri string) (Endpoint, error)

	// Runtime
	NewExchange() *Exchange
}
type DefaultContext struct {
	components map[string]Component
	endpoints  map[string]Endpoint
	mu         sync.RWMutex
	loader     RouteLoader
	routes     []*Route
	started    bool
}

func NewContext() *DefaultContext {
	return &DefaultContext{
		components: make(map[string]Component),
		endpoints:  make(map[string]Endpoint),
		routes:     []*Route{},
	}
}
func (c *DefaultContext) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return nil // idempotent
	}

	// ensure maps are initialized
	if c.components == nil {
		c.components = make(map[string]Component)
	}
	if c.endpoints == nil {
		c.endpoints = make(map[string]Endpoint)
	}

	// Start all routes sequentially. If any route fails to start, return error.
	for _, r := range c.routes {
		if r == nil {
			continue
		}
		if err := r.Start(c); err != nil {
			return fmt.Errorf("failed to start route %s: %w", r.ID, err)
		}
	}

	c.started = true
	return nil
}
func (c *DefaultContext) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		return nil // already stopped
	}

	// Stop routes in reverse order to mirror typical shutdown semantics
	var firstErr error
	for i := len(c.routes) - 1; i >= 0; i-- {
		r := c.routes[i]
		if r == nil {
			continue
		}
		if err := r.Stop(c); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to stop route %s: %w", r.ID, err)
		}
	}

	c.started = false
	return firstErr
}

// SetLoader allows the user to decide how they want to load routes (DSL, YAML, etc.)
func (c *DefaultContext) SetLoader(l RouteLoader) {
	c.loader = l
}

// AddRoutes now uses the assigned loader to ingest definitions.
func (c *DefaultContext) AddRoutes(source interface{}) error {
	// 1. Convert Source (e.g., DSL Builder) into Blueprints (IR)
	definitions, err := c.loader.Load(source)
	if err != nil {
		return fmt.Errorf("loading failed: %w", err)
	}

	for _, def := range definitions {
		// 2. Resolve the Input Endpoint (The "From" part)
		inputEndpoint, err := c.GetEndpoint(def.InputURI)
		if err != nil {
			return err
		}

		// 3. Compile the steps into a chain of Processors
		var pipelineSteps []Processor
		for _, stepDef := range def.Steps {
			// Each definition (To, Choice, etc.) knows how to compile itself
			proc, err := stepDef.Compile(c)
			if err != nil {
				return err
			}
			pipelineSteps = append(pipelineSteps, proc)
		}

		// 4. Wrap steps in a PipelineProcessor
		pipeline := &PipelineProcessor{Children: pipelineSteps}

		// 5. Create the Consumer (The entry point of the route)
		// We pass the pipeline to the consumer so it knows where to send data.
		consumer, err := inputEndpoint.CreateConsumer(pipeline)
		if err != nil {
			return err
		}

		// 6. Finalize the Runtime Route
		runtimeRoute := &Route{
			ID:       def.ID,
			InputURI: def.InputURI,
			Consumer: consumer,
			Pipeline: pipeline,
			context:  c,
		}

		c.routes = append(c.routes, runtimeRoute)

		return nil
	}
	return nil
}

func (c *DefaultContext) RegisterComponent(scheme string, component Component) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.components[scheme] = component
}

func (c *DefaultContext) GetComponent(scheme string) (Component, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	component, exists := c.components[scheme]
	if !exists {
		return nil, errors.New("component not found: " + scheme)
	}
	return component, nil
}

// GetEndpoint resolves a URI into a concrete Endpoint instance.
func (c *DefaultContext) GetEndpoint(uri string) (Endpoint, error) {
	// 1. Thread-safe Cache Lookup
	// We check if this exact URI has been resolved before to ensure we
	// reuse the same Endpoint object (important for resource management).
	c.mu.RLock()
	if ep, ok := c.endpoints[uri]; ok {
		c.mu.RUnlock()
		return ep, nil
	}
	c.mu.RUnlock()

	// 2. Parse the URI to find the Scheme and Path
	// Example: "kafka:my-topic?broker=localhost:9092" -> scheme is "kafka", path is "my-topic"
	scheme, options, err := c.parseUriPathAndOptions(uri)
	if err != nil {
		return nil, err
	}

	// 3. Lookup the Component in the Registry
	component, err := c.GetComponent(scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve endpoint [%s]: %w", uri, err)
	}

	// 4. Delegate Creation to the Component
	// The Component is the expert on its own protocol.
	// options["path"] contains the resource path (e.g., "my-topic", "inbox/file.txt").
	epCfg := EndpointConfig{
		RawURI: uri,
		Scheme: scheme,
		Params: options,
	}
	ep, err := component.CreateEndpoint(epCfg)
	if err != nil {
		return nil, fmt.Errorf("component [%s] could not create endpoint: %w", scheme, err)
	}

	// 5. Store in Cache and Return
	c.mu.Lock()
	c.endpoints[uri] = ep
	c.mu.Unlock()

	return ep, nil
}

// parseUriPathAndOptions handles the logic of extracting scheme, path, and query params.
// It returns scheme, a map of options (including "path"), and any error.
// The "path" key in options contains the resource identifier (e.g., topic name, file path).
// Example: "file:inbox/file.txt?mode=read" -> scheme="file", options={"path": "inbox/file.txt", "mode": "read"}
func (c *DefaultContext) parseUriPathAndOptions(rawUri string) (string, map[string]interface{}, error) {
	// 1. Extract Scheme (e.g., file:, kafka:)
	parts := strings.SplitN(rawUri, ":", 2)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid URI: %s", rawUri)
	}
	scheme := parts[0]

	// 2. Use net/url to parse the "path?query" portion
	u, err := url.Parse(parts[1])
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse URI details: %w", err)
	}

	options := make(map[string]interface{})

	// 3. Extract path (opaque or path component)
	if u.Opaque != "" {
		options["path"] = u.Opaque
	} else if u.Path != "" {
		options["path"] = u.Path
	}

	// 4. Convert url.Values (map[string][]string) to map[string]interface{}
	queryParams := u.Query()
	for k, v := range queryParams {
		if len(v) > 0 {
			options[k] = v[0] // Take the first value for simplicity
		}
	}

	return scheme, options, nil
}

func (c *DefaultContext) compileRoute(def *RouteDefinition) (*Route, error) {
	// Logic to call def.Compile() and wire the Consumer
	return &Route{}, nil
}
func (c *DefaultContext) NewExchange() *Exchange {
	return &Exchange{}
}
