package core

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()

	if msg == nil {
		t.Errorf("expected non-nil Message")
	}

	if msg.body != nil {
		t.Errorf("expected body to be nil initially, got %v", msg.body)
	}

	if msg.headers == nil {
		t.Errorf("expected headers map to be initialized, got nil")
	}

	if len(msg.headers) != 0 {
		t.Errorf("expected empty headers map, got len %d", len(msg.headers))
	}
}

func TestNewMessage_IndependentInstances(t *testing.T) {
	msg1 := NewMessage()
	msg2 := NewMessage()

	msg1.SetHeader("key", "value1")
	msg2.SetHeader("key", "value2")

	if msg1.Header("key") != "value1" {
		t.Errorf("expected msg1 header to be 'value1'")
	}

	if msg2.Header("key") != "value2" {
		t.Errorf("expected msg2 header to be 'value2'")
	}
}

func TestNewMessage_CanSetAndGet(t *testing.T) {
	msg := NewMessage()
	msg.SetBody("test data")
	msg.SetHeader("h1", "v1")

	if msg.Body() != "test data" {
		t.Errorf("expected body 'test data', got %v", msg.Body())
	}

	if msg.Header("h1") != "v1" {
		t.Errorf("expected header 'v1', got %v", msg.Header("h1"))
	}
}

func TestNewExchange(t *testing.T) {
	ex := NewExchange()

	if ex == nil {
		t.Errorf("expected non-nil Exchange")
	}

	if ex.ID() != "1" {
		t.Errorf("expected ID '1', got %s", ex.ID())
	}

	if ex.In() == nil {
		t.Errorf("expected non-nil In message")
	}

	if ex.Out() == nil {
		t.Errorf("expected non-nil Out message")
	}

	if ex.Properties() == nil {
		t.Errorf("expected non-nil properties map")
	}

	if len(ex.Properties()) != 0 {
		t.Errorf("expected empty properties initially")
	}

	if ex.Error() != nil {
		t.Errorf("expected no error initially")
	}
}

func TestNewExchange_InOutAreInitialized(t *testing.T) {
	ex := NewExchange()

	// In and Out should be usable Message instances
	ex.In().SetBody("input")
	ex.In().SetHeader("in_header", "in_value")

	ex.Out().SetBody("output")
	ex.Out().SetHeader("out_header", "out_value")

	if ex.In().Body() != "input" {
		t.Errorf("expected In body 'input'")
	}

	if ex.Out().Body() != "output" {
		t.Errorf("expected Out body 'output'")
	}

	if ex.In().Header("in_header") != "in_value" {
		t.Errorf("expected In header set")
	}

	if ex.Out().Header("out_header") != "out_value" {
		t.Errorf("expected Out header set")
	}
}

func TestNewExchange_IndependentInstances(t *testing.T) {
	ex1 := NewExchange()
	ex2 := NewExchange()

	ex1.In().SetBody("data1")
	ex2.In().SetBody("data2")

	if ex1.In().Body() != "data1" {
		t.Errorf("expected ex1 In body 'data1'")
	}

	if ex2.In().Body() != "data2" {
		t.Errorf("expected ex2 In body 'data2'")
	}

	// In messages should be different instances
	if ex1.In() == ex2.In() {
		t.Errorf("expected independent In messages")
	}
}

func TestMessage_Body(t *testing.T) {
	msg := &Message{
		body:    "test body",
		headers: make(map[string]interface{}),
	}

	if msg.Body() != "test body" {
		t.Errorf("expected 'test body', got %v", msg.Body())
	}
}
func TestMessage_Body_num(t *testing.T) {
	msg := &Message{
		body:    123.45,
		headers: make(map[string]interface{}),
	}

	if msg.Body() != 123.45 {
		t.Errorf("expected number 123.45, got %v", msg.Body())
	}
}

func TestMessage_SetBody(t *testing.T) {
	msg := &Message{headers: make(map[string]interface{})}
	msg.SetBody("new body")

	if msg.Body() != "new body" {
		t.Errorf("expected 'new body', got %v", msg.Body())
	}
}

func TestMessage_Headers(t *testing.T) {
	headers := make(map[string]interface{})
	msg := &Message{
		body:    "test",
		headers: headers,
	}
	if msg.Headers() == nil {
		t.Errorf("expected headers map, got nil")
	}
}

func TestMessage_SetHeader(t *testing.T) {
	msg := &Message{headers: make(map[string]interface{})}
	msg.SetHeader("key1", "value1")
	msg.SetHeader("key2", 42)

	if msg.Header("key1") != "value1" {
		t.Errorf("expected 'value1', got %v", msg.Header("key1"))
	}
	if msg.Header("key2") != 42 {
		t.Errorf("expected 42, got %v", msg.Header("key2"))
	}
}

func TestMessage_Header_NotFound(t *testing.T) {
	msg := &Message{headers: make(map[string]interface{})}
	if msg.Header("nonexistent") != nil {
		t.Errorf("expected nil for nonexistent header")
	}
}

func TestExchange_ID(t *testing.T) {
	ex := &Exchange{id: "exchange-123"}
	if ex.ID() != "exchange-123" {
		t.Errorf("expected 'exchange-123', got %s", ex.ID())
	}
}

func TestExchange_In_Out(t *testing.T) {
	msgIn := &Message{body: "input", headers: make(map[string]interface{})}
	msgOut := &Message{body: "output", headers: make(map[string]interface{})}

	ex := &Exchange{
		id:         "test",
		in:         msgIn,
		out:        msgOut,
		properties: make(map[string]interface{}),
	}
	//msgIn.Body()
	if ex.In().Body() != "input" {
		t.Errorf("expected 'input', got %v", ex.In().Body())
	}
	if ex.Out().Body() != "output" {
		t.Errorf("expected 'output', got %v", ex.Out().Body())
	}
}

func TestExchange_SetIn_SetOut(t *testing.T) {
	ex := &Exchange{
		id:         "test",
		in:         &Message{headers: make(map[string]interface{})},
		out:        &Message{headers: make(map[string]interface{})},
		properties: make(map[string]interface{}),
	}

	newMsgIn := &Message{body: "new input", headers: make(map[string]interface{})}
	newMsgOut := &Message{body: "new output", headers: make(map[string]interface{})}

	ex.SetIn(newMsgIn)
	ex.SetOut(newMsgOut)

	if ex.In().Body() != "new input" {
		t.Errorf("expected 'new input', got %v", ex.In().Body())
	}
	if ex.Out().Body() != "new output" {
		t.Errorf("expected 'new output', got %v", ex.Out().Body())
	}
}

func TestExchange_Properties(t *testing.T) {
	ex := &Exchange{
		id:         "test",
		properties: make(map[string]interface{}),
	}

	ex.SetProperty("prop1", "value1")
	ex.SetProperty("prop2", 100)

	if ex.GetProperty("prop1") != "value1" {
		t.Errorf("expected 'value1', got %v", ex.GetProperty("prop1"))
	}
	if ex.GetProperty("prop2") != 100 {
		t.Errorf("expected 100, got %v", ex.GetProperty("prop2"))
	}
}

func TestExchange_Error(t *testing.T) {
	ex := &Exchange{id: "test", properties: make(map[string]interface{})}

	if ex.Error() != nil {
		t.Errorf("expected no error initially")
	}

	testErr := NewTestError("test error")
	ex.SetError(testErr)

	if ex.Error() != testErr {
		t.Errorf("expected error to be set")
	}
}

func TestExchange_Clone(t *testing.T) {
	msgIn := &Message{body: "input body", headers: map[string]interface{}{"h1": "v1"}}
	msgOut := &Message{body: "output body", headers: map[string]interface{}{"h2": "v2"}}

	original := &Exchange{
		id:         "original",
		in:         msgIn,
		out:        msgOut,
		properties: map[string]interface{}{"p1": "pv1"},
	}

	cloned := original.Clone()

	// Check ID has "_clone" suffix
	if cloned.ID() != "original_clone" {
		t.Errorf("expected 'original_clone', got %s", cloned.ID())
	}

	// Check message bodies are copied
	if cloned.In().Body() != "input body" {
		t.Errorf("expected cloned In to have 'input body'")
	}
	if cloned.Out().Body() != "output body" {
		t.Errorf("expected cloned Out to have 'output body'")
	}

	// Check properties are copied
	if cloned.GetProperty("p1") != "pv1" {
		t.Errorf("expected cloned properties to be copied")
	}

	// Modify original and verify clone is independent
	original.SetProperty("p1", "modified")
	if cloned.GetProperty("p1") == "modified" {
		t.Errorf("clone should not be affected by original modifications")
	}
}

func TestExchange_Clone_IndependentHeaders(t *testing.T) {
	msgIn := &Message{body: "input", headers: map[string]interface{}{"h1": "v1"}}
	msgOut := &Message{body: "output", headers: map[string]interface{}{"h2": "v2"}}

	original := &Exchange{
		id:         "test",
		in:         msgIn,
		out:        msgOut,
		properties: make(map[string]interface{}),
	}

	cloned := original.Clone()

	// Modify original headers
	original.In().SetHeader("h1", "modifie")
	if cloned.In().Header("h1") == "modified" {
		t.Errorf("cloned headers should be independent")
	}
}

// Helper for testing errors
type TestError struct {
	msg string
}

func NewTestError(msg string) *TestError {
	return &TestError{msg: msg}
}

func (e *TestError) Error() string {
	return e.msg
}
