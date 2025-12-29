package file

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sonyjop/camelgo/core"
)

// FileConsumer reads messages from a file and passes them to a processor.
type FileConsumer struct {
	endpoint *FileEndpoint
	target   core.Processor
	file     *os.File
	running  bool
	stopChan chan struct{}
}

func NewFileConsumer(endpoint *FileEndpoint, target core.Processor) *FileConsumer {
	return &FileConsumer{
		endpoint: endpoint,
		target:   target,
		running:  false,
		stopChan: make(chan struct{}),
	}
}

// Start opens the file and begins reading lines, sending each as a message to the target processor.
func (c *FileConsumer) Start(ctx core.Context) error {
	if c.running {
		return nil // idempotent
	}

	// Extract file path from URI
	filePath := c.endpoint.Params["path"].(string)

	// Open file for reading
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	c.file = f
	c.running = true

	// Start reading in a goroutine
	go c.readLoop(ctx)

	return nil
}

// Stop closes the file and stops reading.
func (c *FileConsumer) Stop(ctx core.Context) error {
	if !c.running {
		return nil
	}

	// Signal read loop to stop
	close(c.stopChan)

	if c.file != nil {
		if err := c.file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}

	c.running = false
	return nil
}

// readLoop reads lines from the file and passes each to the target processor.
func (c *FileConsumer) readLoop(ctx core.Context) {
	scanner := bufio.NewScanner(c.file)

	for scanner.Scan() {
		select {
		case <-c.stopChan:
			return
		default:
		}

		line := scanner.Text()

		// Create an exchange with the line as the message body
		exchange := core.NewExchange()
		exchange.In().SetBody(line)

		// Process the exchange
		if err := c.target.Process(ctx, exchange); err != nil {
			// Log error or propagate as needed
			fmt.Printf("error processing line: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner error: %v\n", err)
	}
}
