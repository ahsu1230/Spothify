package monitorrpc

import (
	"../proto"
)

type MonitorInterface interface {
	RegisterServer(args *monitorproto.RegisterArgs, reply *monitorproto.RegisterReply) error
}

type MonitorRPC struct {
	ms MonitorInterface
}

func NewMonitorRPC(ms MonitorInterface) *MonitorRPC {
	return &MonitorRPC{ms}
}

func (mrpc *MonitorRPC) RegisterServer(args *monitorproto.RegisterArgs, reply *monitorproto.RegisterReply) error {
	return mrpc.ms.RegisterServer(args, reply)
}

