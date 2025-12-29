package file

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sonyjop/camelgo/core"
)

func TestCreateEndpoint(t *testing.T) {
	cfg := core.EndpointConfig{
		RawURI: "file:/tmp/f.txt",
		Scheme: "file",
	}

	comp := NewFileComponent()
	ep, err := comp.CreateEndpoint(cfg)
	if err != nil {
		t.Fatalf("CreateEndpoint error: %v", err)
	}
	if ep.GetURI() != cfg.RawURI {
		t.Fatalf("expected URI %s, got %s", cfg.RawURI, ep.GetURI())
	}
}

func TestFileProducer_WriteAndClose(t *testing.T) {
	d := t.TempDir()
	path := filepath.Join(d, "out.txt")

	cfg := core.EndpointConfig{
		RawURI: "file:" + path,
		Scheme: "file",
		Params: map[string]interface{}{"path": path},
	}

	ep := NewFileEndpoint(cfg)
	prod, err := ep.CreateProducer()
	if err != nil {
		t.Fatalf("CreateProducer error: %v", err)
	}

	if err := prod.Start(nil); err != nil {
		t.Fatalf("producer start error: %v", err)
	}

	ex := core.NewExchange()
	ex.In().SetBody("hello-world")
	if err := prod.Process(nil, ex); err != nil {
		prod.Stop(nil)
		t.Fatalf("producer process error: %v", err)
	}

	if err := prod.Stop(nil); err != nil {
		t.Fatalf("producer stop error: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if string(b) != "hello-world\n" {
		t.Fatalf("unexpected file content: %q", string(b))
	}
}

type mockProcessor struct {
	mu     sync.Mutex
	bodies []string
	wg     *sync.WaitGroup
}

func (m *mockProcessor) Process(ctx core.Context, exchange *core.Exchange) error {
	m.mu.Lock()
	m.bodies = append(m.bodies, exchange.In().Body().(string))
	m.mu.Unlock()
	if m.wg != nil {
		m.wg.Done()
	}
	return nil
}

func TestFileConsumer_ReadLines(t *testing.T) {
	d := t.TempDir()
	path := filepath.Join(d, "in.txt")
	lines := []string{"one", "two", "three"}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()

	cfg := core.EndpointConfig{
		RawURI: "file:" + path,
		Scheme: "file",
		Params: map[string]interface{}{"path": path},
	}

	ep := NewFileEndpoint(cfg)

	wg := &sync.WaitGroup{}
	wg.Add(len(lines))
	mp := &mockProcessor{wg: wg}

	cons, err := ep.CreateConsumer(mp)
	if err != nil {
		t.Fatalf("CreateConsumer error: %v", err)
	}

	if err := cons.Start(nil); err != nil {
		t.Fatalf("consumer start error: %v", err)
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
	case <-time.After(2 * time.Second):
		cons.Stop(nil)
		t.Fatalf("timeout waiting for consumer to process lines")
	}

	if err := cons.Stop(nil); err != nil {
		t.Fatalf("consumer stop error: %v", err)
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()
	if len(mp.bodies) != len(lines) {
		t.Fatalf("expected %d processed lines, got %d", len(lines), len(mp.bodies))
	}
	for i := range lines {
		if mp.bodies[i] != lines[i] {
			t.Fatalf("line %d mismatch: expected %q got %q", i, lines[i], mp.bodies[i])
		}
	}
}
