# SSE (Server-Sent Events) 工具库

基于 Kratos 框架的 SSE 流式响应工具库，用于实现服务端向客户端的实时推送。

## 特性

- 专为 Kratos 框架设计，自动处理 Transport Context
- 支持 JSON 自动序列化
- 处理多行数据，确保 SSE 协议规范
- 内置 Buffer Pool，减少内存分配
- 支持心跳保活机制
- 自动屏蔽 Kratos 全局超时，保持长连接

## 安装

该库作为项目内部工具，无需额外安装。

## 快速开始

### 基础用法

```go
package service

import (
    "context"
    
    pb "your-project/api/xxx/v1"
    "your-project/internal/utils/sse"
)

func (s *XXXService) StreamData(ctx context.Context, req *pb.StreamRequest) (*pb.StreamResponse, error) {
    // 1. 创建 SSE Writer
    writer, streamCtx, err := sse.NewWriter(ctx)
    if err != nil {
        return nil, err
    }

    // 2. 使用 streamCtx 进行后续操作（该 context 没有超时限制）
    go func() {
        defer writer.WriteDone()
        
        for {
            select {
            case <-streamCtx.Done():
                return
            default:
                // 发送数据
                writer.WriteEvent(map[string]string{"message": "hello"})
                time.Sleep(time.Second)
            }
        }
    }()

    // 3. 返回 nil 表示响应已由 SSE Writer 处理
    return nil, nil
}
```

### 在 Kratos HTTP 路由中注册

```go
// 在 http.go 中
func NewHTTPServer(c *conf.Server, svc *service.XXXService) *http.Server {
    srv := http.NewServer(
        http.Address(c.Http.Addr),
        http.Timeout(c.Http.Timeout.AsDuration()),
    )

    // 注册路由
    srv.HandleFunc("/api/v1/stream", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        ctx := r.Context()
        _, _ = svc.StreamData(ctx, nil)
    })

    return srv
}
```

## API 文档

### NewWriter

```go
func NewWriter(ctx context.Context) (*Writer, context.Context, error)
```

从 Kratos context 创建 SSE Writer。

**参数：**
- `ctx`: Kratos 请求 context

**返回：**
- `*Writer`: SSE 写入器
- `context.Context`: 无超时限制的 context（用于长连接操作）
- `error`: 错误信息

**可能的错误：**
- `ErrTransportNotFound`: 无法从 context 获取 transport
- `ErrNotHTTP`: 非 HTTP 请求
- `ErrHTTPTransportFailed`: 无法获取 HTTP transport
- `ErrStreamingNotSupported`: 不支持流式响应

---

### WriteEvent

```go
func (s *Writer) WriteEvent(data any) error
```

发送 SSE 事件，数据会自动进行 JSON 序列化。

```go
// 发送结构体
writer.WriteEvent(struct {
    Content string `json:"content"`
}{Content: "hello"})
// 输出: data: {"content":"hello"}\n\n

// 发送 map
writer.WriteEvent(map[string]int{"count": 42})
// 输出: data: {"count":42}\n\n

// 发送字符串（会带引号）
writer.WriteEvent("hello")
// 输出: data: "hello"\n\n
```

---

### WriteEventWithName

```go
func (s *Writer) WriteEventWithName(eventName string, data any) error
```

发送带事件名称的 SSE 事件。

```go
writer.WriteEventWithName("message", map[string]string{"text": "hi"})
// 输出:
// event: message
// data: {"text":"hi"}
//
```

---

### WriteEventWithID

```go
func (s *Writer) WriteEventWithID(id string, data any) error
```

发送带 ID 的 SSE 事件，支持客户端断线重连。

```go
writer.WriteEventWithID("evt-001", map[string]string{"text": "hi"})
// 输出:
// id: evt-001
// data: {"text":"hi"}
//
```

---

### WriteFullEvent

```go
func (s *Writer) WriteFullEvent(event Event) error
```

发送完整的 SSE 事件（支持 id、event、data 组合）。

```go
writer.WriteFullEvent(sse.Event{
    ID:    "evt-001",
    Event: "update",
    Data:  map[string]int{"value": 100},
})
// 输出:
// id: evt-001
// event: update
// data: {"value":100}
//
```

---

### WriteRawEvent

```go
func (s *Writer) WriteRawEvent(data string) error
```

发送原始字符串数据（不进行 JSON 序列化）。

```go
writer.WriteRawEvent("raw text without quotes")
// 输出: data: raw text without quotes\n\n
```

---

### WriteDone

```go
func (s *Writer) WriteDone() error
```

发送结束标记 `[DONE]`，表示流式传输结束。

```go
writer.WriteDone()
// 输出: data: [DONE]\n\n
```

---

### WriteError

```go
func (s *Writer) WriteError(err error) error
```

发送错误事件。

```go
writer.WriteError(errors.New("something went wrong"))
// 输出:
// event: error
// data: {"error":"something went wrong"}
//
```

---

### WriteComment

```go
func (s *Writer) WriteComment(comment string) error
```

发送 SSE 注释（可用作心跳）。

```go
writer.WriteComment("heartbeat")
// 输出: : heartbeat\n\n
```

---

### SetRetry

```go
func (s *Writer) SetRetry(ms int) error
```

设置客户端重连间隔（毫秒）。

```go
writer.SetRetry(3000) // 3 秒后重连
// 输出: retry: 3000\n\n
```

---

### StartHeartbeat

```go
func (s *Writer) StartHeartbeat(interval time.Duration) func()
```

启动定时心跳，返回停止函数。

```go
stop := writer.StartHeartbeat(30 * time.Second)
defer stop()
```

---

### Stream

```go
func (s *Writer) Stream(handler StreamHandler) error
```

通用流式处理方法。

```go
items := []string{"a", "b", "c"}
idx := 0

writer.Stream(func() (any, error) {
    if idx >= len(items) {
        return nil, io.EOF // 返回 EOF 表示结束
    }
    item := items[idx]
    idx++
    return map[string]string{"item": item}, nil
})
```

---

### StreamFunc

```go
func StreamFunc[T any](ctx context.Context, dataCh <-chan T, errCh <-chan error) error
```

简化的流式处理，基于 channel。

**注意：调用方必须关闭 `dataCh` 或取消 `ctx`，否则此函数会永久阻塞。**

```go
dataCh := make(chan Message)
errCh := make(chan error)

go func() {
    defer close(dataCh)
    for _, msg := range messages {
        dataCh <- msg
    }
}()

return sse.StreamFunc(ctx, dataCh, errCh)
```

## 完整示例：AI 聊天流式响应

```go
package service

import (
    "context"
    "io"
    
    pb "your-project/api/admin/v1"
    "your-project/internal/utils/sse"
)

type ChatService struct {
    aiClient AIClient
}

func (s *ChatService) ChatCompletions(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
    // 创建 SSE Writer
    writer, streamCtx, err := sse.NewWriter(ctx)
    if err != nil {
        return nil, err
    }

    // 启动心跳（可选）
    stopHeartbeat := writer.StartHeartbeat(30 * time.Second)
    defer stopHeartbeat()

    // 调用 AI 服务获取流式响应
    stream, err := s.aiClient.CreateChatStream(streamCtx, req.Messages)
    if err != nil {
        return nil, err
    }

    // 流式输出
    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            writer.WriteError(err)
            return nil, nil
        }

        // 发送数据块
        if err := writer.WriteEvent(chunk); err != nil {
            return nil, nil // 客户端断开连接
        }
    }

    // 发送结束标记
    writer.WriteDone()
    
    return nil, nil
}
```

## 前端接收示例

```javascript
const eventSource = new EventSource('/api/v1/chat/stream');

eventSource.onmessage = (event) => {
    if (event.data === '[DONE]') {
        eventSource.close();
        return;
    }
    
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};

eventSource.addEventListener('error', (event) => {
    const data = JSON.parse(event.data);
    console.error('Error:', data.error);
    eventSource.close();
});

eventSource.onerror = (event) => {
    console.error('Connection error');
    eventSource.close();
};
```

## 注意事项

1. **超时处理**：`NewWriter` 返回的 `streamCtx` 已自动移除超时限制，无需担心 Kratos 全局超时中断连接。

2. **CORS**：SSE 库不设置 CORS 头，请在网关或中间件统一处理。

3. **Nginx 配置**：库已自动设置 `X-Accel-Buffering: no`，但建议在 Nginx 配置中也添加：
   ```nginx
   location /api/v1/stream {
       proxy_buffering off;
       proxy_cache off;
       proxy_read_timeout 3600s;
   }
   ```

4. **错误处理**：所有写入方法返回的错误通常表示客户端已断开连接，可以安全地结束处理。

5. **JSON 编码**：`WriteEvent` 会对所有数据进行 JSON 编码（包括字符串），确保前端可以统一使用 `JSON.parse()` 解析。如需发送原始字符串，请使用 `WriteRawEvent`。

## License

Internal use only.

