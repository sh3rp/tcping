package rpc

import (
	"github.com/sh3rp/tcping/rpc"
	context "golang.org/x/net/context"
)

type tcpingdServer struct{}

func (server tcpingdServer) CreateProbe(ctx context.Context, probe *rpc.Probe) (*rpc.Probe, error) {
	panic("not implemented")
}

func (server tcpingdServer) ScheduleProbe(ctx context.Context, schedule *rpc.ProbeSchedule) (*rpc.ProbeSchedule, error) {
	panic("not implemented")
}

func (server tcpingdServer) UnscheduleProbe(ctx context.Context, schedule *rpc.ProbeSchedule) (*rpc.Empty, error) {
	panic("not implemented")
}

func (server tcpingdServer) GetProbeResults(ctx context.Context, query *rpc.ProbeQuery) (*rpc.ProbeQueryResults, error) {
	panic("not implemented")
}

func (server tcpingdServer) StreamProbeResults(probe *rpc.Probe, server rpc.TcpingService_StreamProbeResultsServer) error {
	panic("not implemented")
}
