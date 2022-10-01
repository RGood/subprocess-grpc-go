package runtime

import (
	"context"
	"sync"

	"github.com/RGood/subprocesses-go/internal/protos/runtime"
	"github.com/golang/protobuf/ptypes/empty"
)

type RuntimeService struct {
	runtime.UnimplementedRuntimeServer

	pendingServices sync.Map
}

func NewRuntime() *RuntimeService {
	return &RuntimeService{
		pendingServices: sync.Map{},
	}
}

func (rs *RuntimeService) AddService(id string) {
	serviceChan := make(chan bool)
	rs.pendingServices.Store(id, serviceChan)
}

func (rs *RuntimeService) WaitForService(id string) {
	serviceChan, ok := rs.pendingServices.Load(id)
	if !ok {
		return
	}

	<-serviceChan.(chan bool)
}

func (rs *RuntimeService) Ready(ctx context.Context, msg *runtime.ReadyMessage) (*empty.Empty, error) {
	c, ok := rs.pendingServices.LoadAndDelete(msg.SubprocessId)
	if ok {
		close(c.(chan bool))
	}

	return &empty.Empty{}, nil
}
