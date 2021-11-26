package config

import (
	"context"
	"io/ioutil"

	"github.com/emicklei/xconnect"
	mlog "github.com/microgate-io/microgate/v1/log"
	"gopkg.in/yaml.v2"
)

func Load(fileName string) xconnect.Document {
	ctx := context.Background()
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		mlog.Fatalw(ctx, "failed to read config", "file", fileName, "err", err)
	}
	var doc xconnect.Document
	err = yaml.Unmarshal(content, &doc)
	if err != nil {
		mlog.Fatalw(ctx, "failed to unmarshal config", "err", err)
	}
	return doc
}
