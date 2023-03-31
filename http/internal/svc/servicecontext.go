package svc

import (
	"fmt"
	"github.com/finishy1995/effibot/http/internal/config"
	"github.com/finishy1995/effibot/http/internal/gpt"
	"github.com/finishy1995/effibot/http/internal/scenario"
)

type ServiceContext struct {
	Config config.Config
	GPT    *gpt.Handler
}

func NewServiceContext(c config.Config) *ServiceContext {
	// init scenario
	err := scenario.Init()
	if err != nil {
		panic(err)
	}

	svc := &ServiceContext{
		Config: c,
	}

	svc.GPT = gpt.NewHandler(&c.Spec.GPT)
	if svc.GPT == nil {
		panic(fmt.Errorf("cannot instance GPT"))
	}
	scenario.SetGPTHandler(svc.GPT) // TODO: replace with better way

	return svc
}
