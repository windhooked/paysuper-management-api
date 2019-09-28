package micro

import (
	"context"
	"github.com/Nerufa/go-shared/config"
	"github.com/Nerufa/go-shared/invoker"
	"github.com/Nerufa/go-shared/provider"
	"github.com/google/wire"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, cfg *Config) (*Micro, func(), error) {
	micro := New(ctx, set, cfg)
	return micro, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
		Cfg,
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
	)
)
