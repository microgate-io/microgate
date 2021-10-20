package main

import (
	"context"
	"io/ioutil"

	"github.com/blendle/zapdriver"
	"github.com/emicklei/xconnect"
	"github.com/microgate-io/microgate"
	apilog "github.com/microgate-io/microgate-lib-go/v1/log"
	mlog "github.com/microgate-io/microgate/v1/log"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

func main() {
	// logger, _ := zap.NewProduction()
	cfg := zapdriver.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	defer logger.Sync()
	mlog.InitLogger(logger)

	config := loadConfig()
	apilog.GlobalDebug = config.FindBool("global_debug")

	go microgate.StartInternalProxyServer(config)
	microgate.StartExternalProxyServer(config)
}

func loadConfig() xconnect.Document {
	ctx := context.Background()
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		mlog.Fatalw(ctx, "failed to read config", "err", err)
	}
	var doc xconnect.Document
	err = yaml.Unmarshal(content, &doc)
	if err != nil {
		mlog.Fatalw(ctx, "failed to unmarshal config", "err", err)
	}
	return doc
}
