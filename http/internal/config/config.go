package config

import (
	"github.com/zeromicro/go-zero/rest"
	"time"
)

type Config struct {
	rest.RestConf
	Spec HttpServerConfig
}

type GPTConfig struct {
	Token        string        `json:",optional"`
	TransportUrl string        `json:",optional"`
	Timeout      time.Duration `json:",default=10s"`
	MaxToken     int           `json:",default=4000"`
}

type HttpServerConfig struct {
	GPT GPTConfig
}
