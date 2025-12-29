package main

import (
	"fmt"
	"log"

	"github.com/sonyjop/camelgo/component/file"
	"github.com/sonyjop/camelgo/core"
	"github.com/sonyjop/camelgo/dsl"
)

// MyDSLBuilder defines a custom route using the DSL
type MyDSLBuilder struct {
	*dsl.BaseRouteBuilder
}

func (b *MyDSLBuilder) Configure() {
	// Define a simple route: read from input.txt, write to output.txt
	rd := b.From("file:input.txt")
	dsl.To(rd, "file:output.txt")
}

func main() {
	// Create the context (the runtime engine)
	ctx := &core.DefaultContext{}

	// Register the file component
	fileComp := file.NewFileComponent()
	ctx.RegisterComponent("file", fileComp)

	// Set up the DSL loader
	loader := dsl.NewDSLLoader()
	ctx.SetLoader(loader)

	// Create a route builder and configure routes
	builder := &MyDSLBuilder{
		BaseRouteBuilder: &dsl.BaseRouteBuilder{},
	}

	// Load routes from the builder
	if err := ctx.AddRoutes(builder); err != nil {
		log.Fatalf("failed to add routes: %v", err)
	}

	// Start the context (starts all routes)
	if err := ctx.Start(); err != nil {
		log.Fatalf("failed to start context: %v", err)
	}

	fmt.Println("Routes started successfully")

	// Stop the context (gracefully shuts down all routes)
	if err := ctx.Stop(); err != nil {
		log.Fatalf("failed to stop context: %v", err)
	}

	fmt.Println("Routes stopped successfully")
}
