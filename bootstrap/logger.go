package bootstrap

import (
	conf "fkratos-contrib/api/conf/v1"
	"os"

	aliyunLogger "github.com/go-kratos/kratos/contrib/log/aliyun/v2"
	tencentLogger "github.com/go-kratos/kratos/contrib/log/tencent/v2"
	zapLogger "github.com/go-kratos/kratos/contrib/log/zap/v2"
	zeroLogger "github.com/go-kratos/kratos/contrib/log/zerolog/v2"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

type LoggerType string

const (
	Std     LoggerType = "std"
	Zap     LoggerType = "zap"
	Zerolog LoggerType = "zerolog"
	Aliyun  LoggerType = "aliyun"
	Tencent LoggerType = "tencent"
)

// NewLoggerProvider 创建一个新的日志记录器提供者
func NewLoggerProvider(cfg *conf.Logger, service *Service) log.Logger {
	l := NewLogger(cfg)
	return log.With(
		l,
		"service.id", service.ID,
		"service.name", service.Name,
		"service.version", service.Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}

// NewLogger 创建一个新的日志记录器
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil {
		return NewStdLogger()
	}

	switch LoggerType(cfg.Type) {
	default:
		fallthrough
	case Std:
		return NewStdLogger()
	case Zap:
		return NewZapLogger(cfg)
	case Zerolog:
		return NewZeroLogger()
	case Aliyun:
		return NewAliyunLogger(cfg)
	case Tencent:
		return NewTencentLogger(cfg)
	}
}

// NewStdLogger 创建一个新的日志记录器 - Kratos内置，控制台输出
func NewStdLogger() log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return l
}

// NewZapLogger 创建一个新的日志记录器 - Zap
func NewZapLogger(cfg *conf.Logger) log.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Zap.Filename,
		MaxSize:    int(cfg.Zap.MaxSize),
		MaxBackups: int(cfg.Zap.MaxBackups),
		MaxAge:     int(cfg.Zap.MaxAge),
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	var lvl = new(zapcore.Level)
	if err := lvl.UnmarshalText([]byte(cfg.Zap.Level)); err != nil {
		return nil
	}

	core := zapcore.NewCore(jsonEncoder, writeSyncer, lvl)
	logger := zap.New(core).WithOptions()

	wrapped := zapLogger.NewLogger(logger)

	return wrapped
}

func NewZeroLogger() log.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs // 时间格式
	logger := zerolog.New(os.Stdout)
	return zeroLogger.NewLogger(&logger)
}

// NewAliyunLogger 创建一个新的日志记录器 - Aliyun
func NewAliyunLogger(cfg *conf.Logger) log.Logger {
	wrapped := aliyunLogger.NewAliyunLog(
		aliyunLogger.WithProject(cfg.Aliyun.Project),
		aliyunLogger.WithEndpoint(cfg.Aliyun.Endpoint),
		aliyunLogger.WithAccessKey(cfg.Aliyun.AccessKey),
		aliyunLogger.WithAccessSecret(cfg.Aliyun.AccessSecret),
	)
	return wrapped
}

// NewTencentLogger 创建一个新的日志记录器 - Tencent
func NewTencentLogger(cfg *conf.Logger) log.Logger {
	wrapped, err := tencentLogger.NewLogger(
		tencentLogger.WithTopicID(cfg.Tencent.TopicId),
		tencentLogger.WithEndpoint(cfg.Tencent.Endpoint),
		tencentLogger.WithAccessKey(cfg.Tencent.AccessKey),
		tencentLogger.WithAccessSecret(cfg.Tencent.AccessSecret),
	)
	if err != nil {
		panic(err)
	}
	return wrapped
}
