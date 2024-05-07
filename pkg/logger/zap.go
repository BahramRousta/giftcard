package logger

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.elastic.co/ecszap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"giftcard/config"
	"giftcard/internal/adaptor/logstash"
)

func InitGlobalLogger(lc fx.Lifecycle, L *logstash.LogStash) error {
	// TODO: if you need to write log in file in a critical scenario use FILE and initial a file to pass it as writer in config logger
	//path, err := os.Getwd()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//file := filewriter.NewFile(path+"/logs/"+config.C().Service.Name, os.O_RDWR|os.O_CREATE|os.O_APPEND)
	logger := configLogger(L)
	zap.ReplaceGlobals(logger)

	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			if err := zap.L().Sync(); err != nil {
				log.Println("logger failed to sync:", err)
			}
			log.Println("logger configured successfuly...")
			return nil
		},
	})
	return nil
}

func configLogger(writer io.Writer) *zap.Logger {
	logLevel := getLogLevel()

	wZapCore := ecszap.NewCore(
		ecszap.NewDefaultEncoderConfig(),
		zapcore.AddSync(writer),
		logLevel,
	)
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	terminalZapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		logLevel,
	)
	core := zapcore.NewTee(terminalZapCore, wZapCore)
	logger := zap.New(core, zap.AddCaller())
	return logger.With(zap.String("service.name", config.C().Service.Name))
}

func getLogLevel() zapcore.Level {
	//if config.C().Service {
	//	return zap.DebugLevel
	//}
	return zap.DebugLevel
}

const (
	ZapCtxKey = "zap"
)

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ZapCtxKey).(*zap.Logger)
	if !ok {
		logger = zap.L()
		log.Print("cannot extract logger from context")
	}
	return logger
}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ZapCtxKey, l)
}

func AddToSpan(l *zap.Logger) *otelzap.Logger {
	return otelzap.New(l,
		otelzap.WithMinLevel(zap.DebugLevel),
		otelzap.WithStackTrace(true))
}
