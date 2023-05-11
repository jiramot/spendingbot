package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiramot/spendingbot/env"
	"github.com/jiramot/spendingbot/line"
	"github.com/jiramot/spendingbot/rest"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger := newLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error(err.Error())
		}
	}()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	var r *gin.Engine
	r = gin.Default()

	httpClient := newHTTPClient()
	lineClient, err := linebot.New(env.V.LINE_CHANNEL_SECRET, env.V.LINE_CHANNEL_ACCESS_TOKEN, linebot.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal(err)
	}
	lineSvc := line.NewService(lineClient)
	restHandler := rest.NewRestHandler(lineSvc, env.V.LINE_CHANNEL_SECRET)
	lg := r.Group("/line")
	lg.POST("/webhook", restHandler.HandleWebhook)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:              ":" + env.V.PORT,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 5,
	}

	go func() {
		logger.Info("serve at :" + env.V.PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()

	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")
}

func newLogger() *zap.Logger {
	var config zap.Config
	option := zap.AddStacktrace(zap.PanicLevel)

	switch env.V.LOGGER_MODE {
	case "LOCAL":
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
	default:
		gin.SetMode(gin.ReleaseMode)
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
		}
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build(option)
	if err != nil {
		log.Panic(err)
	}

	return logger
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10000,
			MaxConnsPerHost:     0,
		},
	}
}
