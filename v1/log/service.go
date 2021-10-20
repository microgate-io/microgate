package log

import (
	"context"
	stdlog "log"

	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _logger *zap.Logger = nil

func InitLogger(l *zap.Logger) {
	_logger = l.WithOptions(zap.AddCallerSkip(1)) // because we wrap calls to the logger
	zap.RedirectStdLog(_logger)
}

type LogServiceImpl struct {
	logger *zap.Logger
	apilog.UnimplementedLogServiceServer
}

func NewLogService() *LogServiceImpl {
	if _logger == nil {
		stdlog.Fatal("InitLogger not called")
	}
	return &LogServiceImpl{
		logger: _logger.WithOptions(zap.WithCaller(false)), // caller is passed by request
	}
}

func (s *LogServiceImpl) Log(ctx context.Context, r *apilog.LogRequest) (*emptypb.Empty, error) {
	fields := []zap.Field{}
	for _, each := range r.Attributes {
		if len(each.Key) > 0 {
			if len(each.Value) > 0 {
				fields = append(fields, zap.String(each.Key, each.Value))
			}
		}
	}
	fields = append(fields, zap.String("caller", r.GetCaller()))
	switch r.Level {
	case apilog.LogRequest_INFO:
		s.logger.Info(r.Message, fields...)
	case apilog.LogRequest_WARN:
		s.logger.Warn(r.Message, fields...)
	case apilog.LogRequest_DEBUG:
		s.logger.Debug(r.Message, fields...)
	case apilog.LogRequest_ERROR:
		s.logger.Error(r.Message, fields...)
	}
	return new(emptypb.Empty), nil
}

// TODO remove Sugar calls
func Infow(ctx context.Context, message string, attributes ...interface{}) {
	if _logger == nil {
		stdlog.Println("v1/log/_logger not initialized, call InitLogger", message)
		return
	}
	_logger.Sugar().Infow(message, attributes...)
}

func Warnw(ctx context.Context, message string, attributes ...interface{}) {
	if _logger == nil {
		stdlog.Println("v1/log/_logger not initialized, call InitLogger", message)
		return
	}
	_logger.Sugar().Warnw(message, attributes...)
}

func Debugw(ctx context.Context, message string, attributes ...interface{}) {
	if !apilog.IsDebug(ctx) {
		return
	}
	if _logger == nil {
		stdlog.Println("v1/log/_logger not initialized, call InitLogger", message)
		return
	}
	_logger.Sugar().Debugw(message, attributes...)
}

func Errorw(ctx context.Context, message string, attributes ...interface{}) {
	if _logger == nil {
		stdlog.Println("v1/log/_logger not initialized, call InitLogger", message)
		return
	}
	_logger.Sugar().Errorw(message, attributes...)
}

func Fatalw(ctx context.Context, message string, attributes ...interface{}) {
	if _logger == nil {
		stdlog.Println("v1/log/_logger not initialized, call InitLogger", message)
		return
	}
	_logger.Sugar().Fatalw(message, attributes...)
}
