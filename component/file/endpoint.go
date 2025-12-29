package file

import "github.com/sonyjop/camelgo/core"

type FileEndpoint struct {
	core.EndpointConfig
	//uri       string
	//component *FileComponent
}

func NewFileEndpoint(epCfg core.EndpointConfig) *FileEndpoint {
	return &FileEndpoint{
		EndpointConfig: epCfg,
	}
}
func (e *FileEndpoint) CreateProducer() (core.Producer, error) {
	return NewFileProducer(e), nil
}
func (e *FileEndpoint) CreateConsumer(target core.Processor) (core.Consumer, error) {
	return NewFileConsumer(e, target), nil
}
func (e *FileEndpoint) GetURI() string {
	return e.RawURI
}
