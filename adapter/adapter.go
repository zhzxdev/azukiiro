package adapter

import (
	"context"

	"github.com/zhzxdev/azukiiro/azukiiro/adapter/dummy"
	"github.com/zhzxdev/azukiiro/azukiiro/adapter/uoj"
	"github.com/zhzxdev/azukiiro/azukiiro/common"
)

type JudgeAdapter interface {
	Name() string
	Judge(ctx context.Context, config common.ProblemConfig, problemData string, solutionData string) error
}

var adapters = make(map[string]JudgeAdapter)

func Register(adapter JudgeAdapter) {
	adapters[adapter.Name()] = adapter
}

func Get(name string) (JudgeAdapter, bool) {
	adapter, ok := adapters[name]
	return adapter, ok
}

func init() {
	Register(&dummy.DummyAdapter{})
	Register(&uoj.UojAdapter{})
}
