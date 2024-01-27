package meta

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
)

// GetMetadataFromServer 获取元数据-服务端Context
func GetMetadataFromServer(ctx context.Context, key string) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(key)
	}
	return ""
}

// GetMetadataFromClient 获取元数据-客户端Context
func GetMetadataFromClient(ctx context.Context, key string) string {
	if md, ok := metadata.FromClientContext(ctx); ok {
		return md.Get(key)
	}
	return ""
}

// SetMetadata 设置元数据
func SetMetadata(ctx context.Context, key, value string) context.Context {
	return metadata.AppendToClientContext(ctx, key, value)
}
