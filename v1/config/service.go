package config

import (
	"context"

	apiconfig "github.com/microgate-io/microgate-lib-go/v1/config"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConfigServiceImpl struct {
	apiconfig.UnimplementedConfigServiceServer
}

func NewConfigServiceImpl() *ConfigServiceImpl { return new(ConfigServiceImpl) }

func (s *ConfigServiceImpl) GetConfig(context.Context, *emptypb.Empty) (*apiconfig.GetConfigResponse, error) {
	cfg := &apiconfig.Configuration{
		Entries: make(map[string]*apiconfig.Configuration_Entry),
	}
	cfg.Entries["test-key"] = &apiconfig.Configuration_Entry{StringValue: "test-value"}
	r := &apiconfig.GetConfigResponse{
		Config: cfg,
	}
	return r, nil
}
