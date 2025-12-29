package file

import "github.com/sonyjop/camelgo/core"

type FileComponent struct{}

func NewFileComponent() *FileComponent {
	return &FileComponent{}
}

func (c *FileComponent) GetScheme() string {
	return "file"
}

func (c *FileComponent) CreateEndpoint(epCfg core.EndpointConfig) (core.Endpoint, error) {
	return NewFileEndpoint(epCfg), nil
}
