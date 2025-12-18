package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockResponseWriter 实现 http.ResponseWriter 和 http.Flusher 接口
type mockResponseWriter struct {
	buf        bytes.Buffer
	header     http.Header
	statusCode int
	flushed    int
	mu         sync.Mutex
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		header: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buf.Write(data)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func (m *mockResponseWriter) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushed++
}

func (m *mockResponseWriter) Body() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buf.String()
}

// newTestWriter 创建测试用的 Writer（绕过 Kratos transport）
func newTestWriter() (*Writer, *mockResponseWriter) {
	mock := newMockResponseWriter()
	ctx := context.Background()
	return &Writer{
		w:       mock,
		flusher: mock,
		ctx:     ctx,
	}, mock
}

// =============================================================================
// 测试 sanitizeHeader
// =============================================================================

func TestSanitizeHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "hello", "hello"},
		{"with_newline", "hello\nworld", "hello world"},
		{"with_carriage_return", "hello\rworld", "hello world"},
		{"with_crlf", "hello\r\nworld", "hello  world"},
		{"empty", "", ""},
		{"multiple_newlines", "a\nb\nc", "a b c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeHeader(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeHeader(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 测试 writeData（多行数据处理）
// =============================================================================

func TestWriteData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single_line",
			input:    "hello",
			expected: "data: hello\n",
		},
		{
			name:     "multi_line_lf",
			input:    "line1\nline2",
			expected: "data: line1\ndata: line2\n",
		},
		{
			name:     "multi_line_crlf",
			input:    "line1\r\nline2",
			expected: "data: line1\ndata: line2\n",
		},
		{
			name:     "multi_line_cr",
			input:    "line1\rline2",
			expected: "data: line1\ndata: line2\n",
		},
		{
			name:     "trailing_newline",
			input:    "hello\n",
			expected: "data: hello\ndata: \n",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, _ := newTestWriter()
			buf := &bytes.Buffer{}
			writer.writeData(buf, []byte(tt.input))
			if buf.String() != tt.expected {
				t.Errorf("writeData(%q) = %q, want %q", tt.input, buf.String(), tt.expected)
			}
		})
	}
}

// =============================================================================
// 测试 marshalData
// =============================================================================

func TestMarshalData(t *testing.T) {
	writer, _ := newTestWriter()

	t.Run("string", func(t *testing.T) {
		data, err := writer.marshalData("hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != `"hello"` {
			t.Errorf("marshalData(string) = %q, want %q", string(data), `"hello"`)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type Msg struct {
			Text string `json:"text"`
		}
		data, err := writer.marshalData(Msg{Text: "hi"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != `{"text":"hi"}` {
			t.Errorf("marshalData(struct) = %q, want %q", string(data), `{"text":"hi"}`)
		}
	})

	t.Run("json.RawMessage", func(t *testing.T) {
		raw := json.RawMessage(`{"already":"json"}`)
		data, err := writer.marshalData(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != `{"already":"json"}` {
			t.Errorf("marshalData(RawMessage) = %q, want %q", string(data), `{"already":"json"}`)
		}
	})

	t.Run("map", func(t *testing.T) {
		data, err := writer.marshalData(map[string]int{"count": 42})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != `{"count":42}` {
			t.Errorf("marshalData(map) = %q, want %q", string(data), `{"count":42}`)
		}
	})
}

// =============================================================================
// 测试 WriteEvent
// =============================================================================

func TestWriteEvent(t *testing.T) {
	writer, mock := newTestWriter()

	type Message struct {
		Content string `json:"content"`
	}

	err := writer.WriteEvent(Message{Content: "hello"})
	if err != nil {
		t.Fatalf("WriteEvent failed: %v", err)
	}

	expected := "data: {\"content\":\"hello\"}\n\n"
	if mock.Body() != expected {
		t.Errorf("WriteEvent output = %q, want %q", mock.Body(), expected)
	}

	if mock.flushed != 1 {
		t.Errorf("Flush called %d times, want 1", mock.flushed)
	}
}

// =============================================================================
// 测试 WriteEventWithName
// =============================================================================

func TestWriteEventWithName(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.WriteEventWithName("message", map[string]string{"text": "hi"})
	if err != nil {
		t.Fatalf("WriteEventWithName failed: %v", err)
	}

	body := mock.Body()
	if !strings.Contains(body, "event: message\n") {
		t.Errorf("output should contain 'event: message\\n', got %q", body)
	}
	if !strings.Contains(body, `data: {"text":"hi"}`) {
		t.Errorf("output should contain data, got %q", body)
	}
}

// =============================================================================
// 测试 WriteEventWithID
// =============================================================================

func TestWriteEventWithID(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.WriteEventWithID("123", "test data")
	if err != nil {
		t.Fatalf("WriteEventWithID failed: %v", err)
	}

	body := mock.Body()
	if !strings.Contains(body, "id: 123\n") {
		t.Errorf("output should contain 'id: 123\\n', got %q", body)
	}
}

// =============================================================================
// 测试 WriteFullEvent
// =============================================================================

func TestWriteFullEvent(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.WriteFullEvent(Event{
		ID:    "evt-001",
		Event: "update",
		Data:  map[string]int{"value": 100},
	})
	if err != nil {
		t.Fatalf("WriteFullEvent failed: %v", err)
	}

	body := mock.Body()
	if !strings.Contains(body, "id: evt-001\n") {
		t.Errorf("output should contain id, got %q", body)
	}
	if !strings.Contains(body, "event: update\n") {
		t.Errorf("output should contain event, got %q", body)
	}
	if !strings.Contains(body, `data: {"value":100}`) {
		t.Errorf("output should contain data, got %q", body)
	}
}

// =============================================================================
// 测试 WriteRawEvent
// =============================================================================

func TestWriteRawEvent(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.WriteRawEvent("raw string without json encoding")
	if err != nil {
		t.Fatalf("WriteRawEvent failed: %v", err)
	}

	expected := "data: raw string without json encoding\n\n"
	if mock.Body() != expected {
		t.Errorf("WriteRawEvent output = %q, want %q", mock.Body(), expected)
	}
}

// =============================================================================
// 测试 WriteDone
// =============================================================================

func TestWriteDone(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.WriteDone()
	if err != nil {
		t.Fatalf("WriteDone failed: %v", err)
	}

	expected := "data: [DONE]\n\n"
	if mock.Body() != expected {
		t.Errorf("WriteDone output = %q, want %q", mock.Body(), expected)
	}
}

// =============================================================================
// 测试 WriteError
// =============================================================================

func TestWriteError(t *testing.T) {
	t.Run("with_error", func(t *testing.T) {
		writer, mock := newTestWriter()

		err := writer.WriteError(errors.New("something went wrong"))
		if err != nil {
			t.Fatalf("WriteError failed: %v", err)
		}

		body := mock.Body()
		if !strings.Contains(body, "event: error\n") {
			t.Errorf("output should contain 'event: error', got %q", body)
		}
		if !strings.Contains(body, "something went wrong") {
			t.Errorf("output should contain error message, got %q", body)
		}
	})

	t.Run("nil_error", func(t *testing.T) {
		writer, mock := newTestWriter()

		err := writer.WriteError(nil)
		if err != nil {
			t.Fatalf("WriteError(nil) failed: %v", err)
		}

		if mock.Body() != "" {
			t.Errorf("WriteError(nil) should write nothing, got %q", mock.Body())
		}
	})
}

// =============================================================================
// 测试 WriteComment
// =============================================================================

func TestWriteComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "simple",
			comment:  "heartbeat",
			expected: ": heartbeat\n\n",
		},
		{
			name:     "multiline",
			comment:  "line1\nline2",
			expected: ": line1\n: line2\n\n",
		},
		{
			name:     "with_crlf",
			comment:  "line1\r\nline2",
			expected: ": line1\n: line2\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, mock := newTestWriter()
			err := writer.WriteComment(tt.comment)
			if err != nil {
				t.Fatalf("WriteComment failed: %v", err)
			}
			if mock.Body() != tt.expected {
				t.Errorf("WriteComment(%q) = %q, want %q", tt.comment, mock.Body(), tt.expected)
			}
		})
	}
}

// =============================================================================
// 测试 SetRetry
// =============================================================================

func TestSetRetry(t *testing.T) {
	writer, mock := newTestWriter()

	err := writer.SetRetry(3000)
	if err != nil {
		t.Fatalf("SetRetry failed: %v", err)
	}

	expected := "retry: 3000\n\n"
	if mock.Body() != expected {
		t.Errorf("SetRetry output = %q, want %q", mock.Body(), expected)
	}
}

// =============================================================================
// 测试 Stream
// =============================================================================

func TestStream(t *testing.T) {
	writer, mock := newTestWriter()

	items := []string{"first", "second", "third"}
	idx := 0

	err := writer.Stream(func() (any, error) {
		if idx >= len(items) {
			return nil, io.EOF
		}
		item := items[idx]
		idx++
		return map[string]string{"msg": item}, nil
	})

	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	body := mock.Body()
	for _, item := range items {
		if !strings.Contains(body, item) {
			t.Errorf("output should contain %q, got %q", item, body)
		}
	}
	if !strings.Contains(body, "[DONE]") {
		t.Errorf("output should contain [DONE], got %q", body)
	}
}

// =============================================================================
// 测试 StartHeartbeat
// =============================================================================

func TestStartHeartbeat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		writer, mock := newTestWriter()

		stop := writer.StartHeartbeat(50 * time.Millisecond)
		time.Sleep(130 * time.Millisecond)
		stop()

		body := mock.Body()
		count := strings.Count(body, ": heartbeat\n")
		if count < 2 {
			t.Errorf("expected at least 2 heartbeats, got %d", count)
		}
	})

	t.Run("zero_interval", func(t *testing.T) {
		writer, _ := newTestWriter()

		stop := writer.StartHeartbeat(0)
		stop() // should not panic
	})

	t.Run("negative_interval", func(t *testing.T) {
		writer, _ := newTestWriter()

		stop := writer.StartHeartbeat(-1 * time.Second)
		stop() // should not panic
	})
}

// =============================================================================
// 测试 Buffer Pool
// =============================================================================

func TestBufferPool(t *testing.T) {
	t.Run("normal_size", func(t *testing.T) {
		buf := getBuffer()
		buf.WriteString("test data")
		putBuffer(buf)
		// 正常大小的 buffer 应该被放回池中
	})

	t.Run("large_buffer_not_returned", func(t *testing.T) {
		buf := getBuffer()
		// 写入超过 maxBufferSize 的数据
		largeData := make([]byte, maxBufferSize+1000)
		buf.Write(largeData)
		putBuffer(buf) // 大 buffer 不应该被放回池中
		// 这里主要是验证不会 panic
	})
}

// =============================================================================
// 测试 streamContext
// =============================================================================

func TestStreamContext(t *testing.T) {
	t.Run("no_deadline", func(t *testing.T) {
		// 即使 parent 有超时，streamContext 也不应该有 deadline
		parent, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		ctx := &streamContext{
			values: parent,
			done:   make(chan struct{}),
		}

		deadline, ok := ctx.Deadline()
		if ok {
			t.Errorf("streamContext should return ok=false, got deadline=%v", deadline)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		ctx := &streamContext{
			values: context.Background(),
			done:   make(chan struct{}),
		}

		// 取消 context
		ctx.cancel(context.Canceled)

		select {
		case <-ctx.Done():
			// 期望的行为
		case <-time.After(time.Second):
			t.Error("cancel should close done channel")
		}

		// 验证错误
		if ctx.Err() != context.Canceled {
			t.Errorf("Err() = %v, want %v", ctx.Err(), context.Canceled)
		}
	})

	t.Run("cancel_only_once", func(t *testing.T) {
		ctx := &streamContext{
			values: context.Background(),
			done:   make(chan struct{}),
		}

		// 多次取消不应该 panic
		ctx.cancel(context.Canceled)
		ctx.cancel(context.DeadlineExceeded) // 第二次取消应该被忽略

		// 错误应该是第一次设置的
		if ctx.Err() != context.Canceled {
			t.Errorf("Err() = %v, want %v", ctx.Err(), context.Canceled)
		}
	})

	t.Run("value_propagation", func(t *testing.T) {
		type key string
		parent := context.WithValue(context.Background(), key("test"), "value")

		ctx := &streamContext{
			values: parent,
			done:   make(chan struct{}),
		}

		if ctx.Value(key("test")) != "value" {
			t.Errorf("Value() should propagate from parent context")
		}
	})
}

// =============================================================================
// 测试 Write（io.Writer 接口）
// =============================================================================

func TestWrite(t *testing.T) {
	writer, mock := newTestWriter()

	n, err := writer.Write([]byte("raw bytes"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 9 {
		t.Errorf("Write returned n=%d, want 9", n)
	}
	if mock.Body() != "raw bytes" {
		t.Errorf("Write output = %q, want %q", mock.Body(), "raw bytes")
	}
	if mock.flushed != 1 {
		t.Errorf("Flush called %d times, want 1", mock.flushed)
	}
}

// =============================================================================
// 测试 Context 方法
// =============================================================================

func TestContext(t *testing.T) {
	writer, _ := newTestWriter()

	ctx := writer.Context()
	if ctx == nil {
		t.Error("Context() should not return nil")
	}
}

// =============================================================================
// 集成测试：模拟真实 HTTP 请求
// =============================================================================

func TestHTTPIntegration(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("ResponseWriter does not support Flusher")
		}

		writer := &Writer{
			w:       w,
			flusher: flusher,
			ctx:     r.Context(),
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		_ = writer.WriteEvent(map[string]string{"msg": "hello"})
		_ = writer.WriteDone()
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, `{"msg":"hello"}`) {
		t.Errorf("response should contain message, got %q", bodyStr)
	}
	if !strings.Contains(bodyStr, "[DONE]") {
		t.Errorf("response should contain [DONE], got %q", bodyStr)
	}
}

// =============================================================================
// 并发安全测试
// =============================================================================

func TestConcurrentWrites(t *testing.T) {
	writer, _ := newTestWriter()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = writer.WriteEvent(map[string]int{"n": n})
		}(i)
	}
	wg.Wait()
	// 主要验证并发写入不会 panic
}
