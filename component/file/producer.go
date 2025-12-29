package file

import (
	"fmt"
	"os"

	"github.com/sonyjop/camelgo/core"
)

// FileProducer writes messages to a file.
type FileProducer struct {
	endpoint *FileEndpoint
	file     *os.File
	running  bool
}

func NewFileProducer(endpoint *FileEndpoint) *FileProducer {
	return &FileProducer{
		endpoint: endpoint,
		running:  false,
	}
}

// Start opens the file for writing.
func (p *FileProducer) Start(ctx core.Context) error {
	if p.running {
		return nil // idempotent
	}

	// Extract file path from URI (e.g., "file:/path/to/file.txt" -> "/path/to/file.txt")
	filePath := p.endpoint.Params["path"].(string)

	// Open file in append mode, create if not exists
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	p.file = f
	p.running = true
	return nil
}

// Stop closes the file.
func (p *FileProducer) Stop(ctx core.Context) error {
	if !p.running {
		return nil
	}

	if p.file != nil {
		if err := p.file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}

	p.running = false
	return nil
}

// Process writes the message body to the file.
func (p *FileProducer) Process(ctx core.Context, exchange *core.Exchange) error {
	if !p.running {
		return fmt.Errorf("FileProducer not started")
	}

	body := exchange.In().Body()
	if body == nil {
		return nil
	}

	// Convert body to string
	content := fmt.Sprintf("%v\n", body)

	_, err := p.file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

/*// extractFilePath extracts the file path from a URI like "file:/path/to/file.txt"
func (p *FileProducer) extractFilePath(uri string) string {
	// Simple extraction: remove "file:" prefix
	if len(uri) > 5 && uri[:5] == "file:" {
		return uri[5:]
	}
	return uri
}*/
