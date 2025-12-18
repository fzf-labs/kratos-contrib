// Package sse 提供基于 Kratos 框架的 SSE（Server-Sent Events）流式响应工具
package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// 预定义错误
var (
	ErrTransportNotFound     = errors.New("sse: failed to get transport from context")
	ErrNotHTTP               = errors.New("sse: only HTTP requests are supported")
	ErrHTTPTransportFailed   = errors.New("sse: failed to get http transport")
	ErrStreamingNotSupported = errors.New("sse: streaming not supported")
)

// 预定义字节切片，避免重复分配
var (
	dataPrefix    = []byte("data: ")
	eventPrefix   = []byte("event: ")
	idPrefix      = []byte("id: ")
	retryPrefix   = []byte("retry: ")
	commentPrefix = []byte(": ")
	lineEnd       = []byte("\n")
	eventEnd      = []byte("\n\n")
	doneMessage   = []byte("data: [DONE]\n\n")
)

// buffer 池，用于减少内存分配
var bufferPool = sync.Pool{
	New: func() any {
		b := new(bytes.Buffer)
		b.Grow(256) // 预分配 256 字节，适合大多数 SSE 消息
		return b
	},
}

// getBuffer 从池中获取 buffer
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// maxBufferSize 限制放回池中的 buffer 最大容量，防止大 buffer 长期占用内存
const maxBufferSize = 4096

// putBuffer 将 buffer 放回池中
func putBuffer(buf *bytes.Buffer) {
	// 如果 buffer 容量过大，丢弃它让 GC 回收
	if buf.Cap() > maxBufferSize {
		return
	}
	bufferPool.Put(buf)
}

// Writer SSE 流式写入器
type Writer struct {
	mu      sync.Mutex
	w       stdhttp.ResponseWriter
	flusher stdhttp.Flusher
	ctx     context.Context
}

// noDeadlineContext 屏蔽 Deadline 但保留 Cancel
type noDeadlineContext struct {
	context.Context
}

func (c *noDeadlineContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// NewWriter 从 Kratos context 中创建 SSE Writer
// 返回 Writer 和一个无超时的 context（用于长时间流式操作）
func NewWriter(ctx context.Context) (*Writer, context.Context, error) {
	// 使用 transport.FromServerContext 获取传输层信息
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, nil, ErrTransportNotFound
	}

	// 判断是否为 HTTP 请求
	if tr.Kind() != transport.KindHTTP {
		return nil, nil, ErrNotHTTP
	}

	// 获取 HTTP Transport
	httpTransport, ok := tr.(*http.Transport)
	if !ok {
		return nil, nil, ErrHTTPTransportFailed
	}

	// 获取 ResponseWriter
	w := httpTransport.Response()

	// 先检查 Flusher 支持，避免写入响应头后才发现不支持流式
	flusher, ok := w.(stdhttp.Flusher)
	if !ok {
		return nil, nil, ErrStreamingNotSupported
	}

	// 设置 SSE 响应头（CORS 由网关/中间件统一处理）
	header := httpTransport.ReplyHeader()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no") // 禁用 Nginx 缓存，确保流式实时性

	// 写入 200 状态码（在确认支持流式后再写）
	w.WriteHeader(stdhttp.StatusOK)

	// 关键：保留取消信号，但移除超时限制
	// 这样当客户端关闭连接时，ctx.Done() 依然能触发，但不会被 Kratos 的全局 Timeout 中断
	streamCtx := &noDeadlineContext{Context: ctx}

	return &Writer{
		w:       w,
		flusher: flusher,
		ctx:     streamCtx,
	}, streamCtx, nil
}

// Context 返回无超时的流式 context
func (s *Writer) Context() context.Context {
	return s.ctx
}

// writeAndFlush 写入数据并刷新
func (s *Writer) writeAndFlush(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.w.Write(data); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

// Write 写入原始字节数据并立即刷新（实现 io.Writer 接口）
func (s *Writer) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, err := s.w.Write(data)
	if err != nil {
		return n, err
	}
	s.flusher.Flush()
	return n, nil
}

// marshalData 格式化数据为 JSON
// 注意：所有类型都会进行 JSON 编码，确保前端可以统一使用 JSON.parse() 解析
// 如需发送原始字符串，请使用 WriteRawEvent 方法
func (s *Writer) marshalData(data any) ([]byte, error) {
	// 如果是 json.RawMessage，直接返回（已经是有效的 JSON）
	if raw, ok := data.(json.RawMessage); ok {
		return raw, nil
	}
	return json.Marshal(data)
}

// sanitizeHeader 清理头部字段中的换行符，防止破坏 SSE 协议结构
func sanitizeHeader(s string) string {
	return strings.NewReplacer("\n", " ", "\r", " ").Replace(s)
}

// writeData 内部方法：处理多行数据，确保每一行都带有 data: 前缀
// 兼容 \n, \r, \r\n 换行符
func (s *Writer) writeData(buf *bytes.Buffer, data []byte) {
	for len(data) > 0 {
		buf.Write(dataPrefix)
		idx := bytes.IndexAny(data, "\n\r")
		if idx == -1 {
			buf.Write(data)
			buf.Write(lineEnd)
			break
		}
		buf.Write(data[:idx])
		buf.Write(lineEnd)

		// 处理 \r\n 序列
		if data[idx] == '\r' && idx+1 < len(data) && data[idx+1] == '\n' {
			data = data[idx+2:]
		} else {
			data = data[idx+1:]
		}

		if len(data) == 0 {
			// 如果数据以换行符结尾，按照规范补一个空 data 行
			buf.Write(dataPrefix)
			buf.Write(lineEnd)
		}
	}
}

// WriteEvent 写入 SSE 格式的事件数据
// data: 要发送的数据
func (s *Writer) WriteEvent(data any) error {
	jsonData, err := s.marshalData(data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %w", err)
	}

	buf := getBuffer()
	defer putBuffer(buf)

	s.writeData(buf, jsonData)
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// WriteEventWithName 写入带事件名称的 SSE 格式数据
// eventName: 事件名称
// data: 要发送的数据
func (s *Writer) WriteEventWithName(eventName string, data any) error {
	jsonData, err := s.marshalData(data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %w", err)
	}

	buf := getBuffer()
	defer putBuffer(buf)

	buf.Write(eventPrefix)
	buf.WriteString(sanitizeHeader(eventName))
	buf.Write(lineEnd)
	s.writeData(buf, jsonData)
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// WriteEventWithID 写入带 ID 的 SSE 事件（支持断线重连）
func (s *Writer) WriteEventWithID(id string, data any) error {
	jsonData, err := s.marshalData(data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %w", err)
	}

	buf := getBuffer()
	defer putBuffer(buf)

	buf.Write(idPrefix)
	buf.WriteString(sanitizeHeader(id))
	buf.Write(lineEnd)
	s.writeData(buf, jsonData)
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// Event SSE 完整事件结构
type Event struct {
	ID    string // 事件 ID（可选）
	Event string // 事件名称（可选）
	Data  any    // 事件数据
}

// WriteFullEvent 写入完整的 SSE 事件（支持 id、event、data 组合）
func (s *Writer) WriteFullEvent(event Event) error {
	jsonData, err := s.marshalData(event.Data)
	if err != nil {
		return fmt.Errorf("marshal data failed: %w", err)
	}

	buf := getBuffer()
	defer putBuffer(buf)

	if event.ID != "" {
		buf.Write(idPrefix)
		buf.WriteString(sanitizeHeader(event.ID))
		buf.Write(lineEnd)
	}
	if event.Event != "" {
		buf.Write(eventPrefix)
		buf.WriteString(sanitizeHeader(event.Event))
		buf.Write(lineEnd)
	}
	s.writeData(buf, jsonData)
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// WriteRawEvent 写入原始字符串数据（不进行 JSON 序列化）
func (s *Writer) WriteRawEvent(data string) error {
	buf := getBuffer()
	defer putBuffer(buf)

	s.writeData(buf, []byte(data))
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// WriteDone 发送结束标记 [DONE]
func (s *Writer) WriteDone() error {
	return s.writeAndFlush(doneMessage)
}

// WriteError 发送错误信息
func (s *Writer) WriteError(err error) error {
	if err == nil {
		return nil
	}
	errData := map[string]string{"error": err.Error()}
	return s.WriteEventWithName("error", errData)
}

// WriteComment 发送 SSE 注释（可用作心跳）
// 正确处理 \n, \r, \r\n 换行符
func (s *Writer) WriteComment(comment string) error {
	buf := getBuffer()
	defer putBuffer(buf)

	data := []byte(comment)
	for len(data) > 0 {
		buf.Write(commentPrefix)
		idx := bytes.IndexAny(data, "\n\r")
		if idx == -1 {
			buf.Write(data)
			buf.Write(lineEnd)
			break
		}
		buf.Write(data[:idx])
		buf.Write(lineEnd)

		// 处理 \r\n 序列
		if data[idx] == '\r' && idx+1 < len(data) && data[idx+1] == '\n' {
			data = data[idx+2:]
		} else {
			data = data[idx+1:]
		}
	}
	buf.Write(lineEnd)

	return s.writeAndFlush(buf.Bytes())
}

// SetRetry 设置客户端重连间隔（毫秒）
func (s *Writer) SetRetry(ms int) error {
	buf := getBuffer()
	defer putBuffer(buf)

	buf.Write(retryPrefix)
	buf.WriteString(strconv.Itoa(ms))
	buf.Write(eventEnd)

	return s.writeAndFlush(buf.Bytes())
}

// StartHeartbeat 启动心跳（返回停止函数）
func (s *Writer) StartHeartbeat(interval time.Duration) func() {
	if interval <= 0 {
		return func() {}
	}
	ticker := time.NewTicker(interval)
	done := make(chan struct{})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.WriteComment("heartbeat"); err != nil {
					// 如果写入失败（如连接已断开），停止心跳
					return
				}
			case <-done:
				return
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return func() {
		select {
		case <-done:
		default:
			close(done)
		}
	}
}

// StreamHandler 流式数据处理函数类型
// 返回要发送的数据和错误
// 当返回 io.EOF 时表示流结束
type StreamHandler func() (data any, err error)

// Stream 通用流式处理方法
// handler: 每次调用返回一个数据块，返回 io.EOF 表示结束
func (s *Writer) Stream(handler StreamHandler) error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			data, err := handler()
			if err == io.EOF {
				return s.WriteDone()
			}
			if err != nil {
				return err
			}
			if err := s.WriteEvent(data); err != nil {
				return err
			}
		}
	}
}

// StreamFunc 简化的流式处理，接收一个 channel
// 注意：调用方必须关闭 dataCh 或取消 ctx，否则此函数会永久阻塞
func StreamFunc[T any](ctx context.Context, dataCh <-chan T, errCh <-chan error) error {
	writer, streamCtx, err := NewWriter(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case data, ok := <-dataCh:
			if !ok {
				return writer.WriteDone()
			}
			if err := writer.WriteEvent(data); err != nil {
				return err
			}
		case err, ok := <-errCh:
			if !ok {
				// 如果 errCh 关闭，设为 nil 避免 select 频繁触发导致 CPU 飙升
				errCh = nil
				continue
			}
			if err != nil {
				return err
			}
		case <-streamCtx.Done():
			return streamCtx.Err()
		}
	}
}
