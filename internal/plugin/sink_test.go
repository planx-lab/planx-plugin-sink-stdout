package plugin

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

// captureStdout swaps os.Stdout for a pipe, runs fn, and returns whatever was
// written. Lets the sink be tested without a real downstream consumer.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = orig
	return <-done
}

func TestSink_New_ReturnsSinkSPI(t *testing.T) {
	var _ sdk.SinkSPI = New()
}

func TestSink_Init_NoConfig(t *testing.T) {
	if err := (&Sink{}).Init(context.Background(), nil); err != nil {
		t.Fatalf("Init: %v", err)
	}
}

func TestSink_WriteBatch_PrettyPrintsRows(t *testing.T) {
	s := &Sink{}
	out := captureStdout(t, func() {
		if err := s.WriteBatch(sdk.Rows{
			{"name": "alice", "id": "1"},
			{"name": "bob", "id": "2"},
		}); err != nil {
			t.Fatalf("WriteBatch: %v", err)
		}
	})
	// one [SINK] row line per input row; keys sorted alphabetically.
	if !strings.Contains(out, "[SINK] row: id=1 name=alice") {
		t.Errorf("missing alice row in output:\n%s", out)
	}
	if !strings.Contains(out, "[SINK] row: id=2 name=bob") {
		t.Errorf("missing bob row in output:\n%s", out)
	}
	if !strings.Contains(out, "[SINK] row:") {
		t.Errorf("output should use the row prefix:\n%s", out)
	}
}

func TestSink_WriteBatch_EmptyRows_NoOutput(t *testing.T) {
	s := &Sink{}
	out := captureStdout(t, func() {
		if err := s.WriteBatch(sdk.Rows{}); err != nil {
			t.Fatalf("WriteBatch: %v", err)
		}
	})
	if out != "" {
		t.Errorf("empty Rows should print nothing, got %q", out)
	}
}

func TestSink_WriteBatch_DBBatchFallback(t *testing.T) {
	s := &Sink{}
	out := captureStdout(t, func() {
		if err := s.WriteBatch(DBBatch{
			Columns: []string{"id", "name"},
			Rows:    []DBRow{{Types: []byte{3, 3}, Vals: []string{"1", "alice"}}},
		}); err != nil {
			t.Fatalf("WriteBatch: %v", err)
		}
	})
	if !strings.Contains(out, "id=1 name=alice") {
		t.Errorf("DBBatch fallback should pretty-print column=value:\n%s", out)
	}
}

func TestSink_WriteBatch_GenericFallback(t *testing.T) {
	s := &Sink{}
	out := captureStdout(t, func() {
		if err := s.WriteBatch(42); err != nil {
			t.Fatalf("WriteBatch: %v", err)
		}
	})
	if !strings.Contains(out, "[SINK] Received Batch: 42") {
		t.Errorf("generic fallback should %%v-dump:\n%s", out)
	}
}

func TestSink_Close_NoError(t *testing.T) {
	if err := (&Sink{}).Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
