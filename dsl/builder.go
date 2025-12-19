package dsl

import "github.com/sonyjop/camelgo/core"

// BaseRouteBuilder provides common logic for user-defined routes.
type BaseRouteBuilder struct {
	definitions []*core.RouteDefinition
}

func (b *BaseRouteBuilder) From(uri string) *core.RouteDefinition {
	route := &core.RouteDefinition{InputURI: uri}
	b.definitions = append(b.definitions, route)
	return route
}

func (b *BaseRouteBuilder) GetRouteDefinitions() []*core.RouteDefinition {
	return b.definitions
}
