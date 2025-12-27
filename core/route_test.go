package core

import (
	"testing"
)

// MockProcessor is a test double for Processor
type MockProcessor struct {
	ProcessCalled bool
	ProcessErr    error
	ProcessCount  int
}

func (m *MockProcessor) Process(ctx Context, exchange *Exchange) error {
	m.ProcessCalled = true
	m.ProcessCount++
	return m.ProcessErr
}

func TestRoute_Start_NilConsumer(t *testing.T) {
	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: nil,
	}

	err := route.Start(nil)
	if err != nil {
		t.Errorf("expected no error for nil consumer, got %v", err)
	}
}

func TestRoute_Stop_NilConsumer(t *testing.T) {
	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: nil,
	}

	err := route.Stop(nil)
	if err != nil {
		t.Errorf("expected no error for nil consumer, got %v", err)
	}
}

func TestRoute_Start_WithConsumer(t *testing.T) {
	mockConsumer := &MockConsumer{}

	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: mockConsumer,
	}

	err := route.Start(nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !mockConsumer.StartCalled {
		t.Errorf("expected Consumer.Start() to be called")
	}
}

func TestRoute_Stop_WithConsumer(t *testing.T) {
	mockConsumer := &MockConsumer{}

	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: mockConsumer,
	}

	err := route.Stop(nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !mockConsumer.StopCalled {
		t.Errorf("expected Consumer.Stop() to be called")
	}
}

func TestRoute_Start_ConsumerError(t *testing.T) {
	mockConsumer := &MockConsumer{
		StartErr: NewTestError("start failed"),
	}

	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: mockConsumer,
	}

	err := route.Start(nil)
	if err == nil {
		t.Errorf("expected error from Consumer.Start()")
	}

	if err.Error() != "start failed" {
		t.Errorf("expected 'start failed', got %v", err)
	}
}

func TestRoute_Stop_ConsumerError(t *testing.T) {
	mockConsumer := &MockConsumer{
		StopErr: NewTestError("stop failed"),
	}

	route := &Route{
		ID:       "test-route",
		InputURI: "test://input",
		Consumer: mockConsumer,
	}

	err := route.Stop(nil)
	if err == nil {
		t.Errorf("expected error from Consumer.Stop()")
	}

	if err.Error() != "stop failed" {
		t.Errorf("expected 'stop failed', got %v", err)
	}
}

func TestRoute_Attributes(t *testing.T) {
	mockConsumer := &MockConsumer{}
	mockProcessor := &MockProcessor{}

	route := &Route{
		ID:       "route-123",
		InputURI: "kafka://my-topic",
		Consumer: mockConsumer,
		Pipeline: mockProcessor,
	}

	if route.ID != "route-123" {
		t.Errorf("expected ID 'route-123', got %s", route.ID)
	}

	if route.InputURI != "kafka://my-topic" {
		t.Errorf("expected InputURI 'kafka://my-topic', got %s", route.InputURI)
	}
}

// MockConsumer for testing
type MockConsumer struct {
	StartCalled bool
	StopCalled  bool
	StartErr    error
	StopErr     error
}

func (m *MockConsumer) Start(ctx Context) error {
	m.StartCalled = true
	return m.StartErr
}

func (m *MockConsumer) Stop(ctx Context) error {
	m.StopCalled = true
	return m.StopErr
}
