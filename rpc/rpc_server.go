package rpc

import (
	"sync"

	"github.com/google/uuid"
	"github.com/sh3rp/tcping/tcping"
	context "golang.org/x/net/context"
	cron "gopkg.in/robfig/cron.v2"
)

type TcpingdServer struct {
	probes      map[string]tcping.Probe
	scheduler   *cron.Cron
	probesLock  sync.Mutex
	resultQueue chan probeResult
}

type probeResult struct {
	result tcping.ProbeResult
	err    error
}

func (server TcpingdServer) CreateProbe(ctx context.Context, probe *Probe) (*Probe, error) {
	if probe.Id == "" {
		probe.Id = uuid.New().String()
	}

	server.probesLock.Lock()
	if _, ok := server.probes[probe.Id]; !ok {
		server.probes[probe.Id] = tcping.NewProbe("0.0.0.0", probe.Host, 3000, uint16(probe.Port), false)
	}
	server.probesLock.Unlock()

	panic("not implemented")
}

func (server TcpingdServer) ScheduleProbe(ctx context.Context, schedule *ProbeSchedule) (*ProbeSchedule, error) {
	if schedule.Schedule != "" {
		// check if schedule id was supplied and modify it if so
		// otherwise add the new schedule
		if schedule.Id == "" {
			probe, probeExists := server.probes[schedule.Probe.Id]
			if schedule.Probe != nil && schedule.Probe.Id != "" && probeExists {
				entryId, err := server.scheduler.AddFunc(schedule.Schedule, func() {
					result, err := probe.GetLatency()
					server.resultQueue <- probeResult{result, err}
				})
				if err != nil {
					return nil, err
				}
				schedule.Id = string(entryId)
			} else {

			}
		} else {
			_, probeExists := server.probes[schedule.Probe.Id]
			if schedule.Probe != nil && schedule.Probe.Id != "" && probeExists {
				// delete schedule then reschedule
			} else {

			}
		}
	}
	return schedule, nil
}

func (server TcpingdServer) UnscheduleProbe(ctx context.Context, schedule *ProbeSchedule) (*Empty, error) {
	panic("not implemented")
}

func (server TcpingdServer) GetProbeResults(ctx context.Context, query *ProbeQuery) (*ProbeQueryResults, error) {
	panic("not implemented")
}

func (server TcpingdServer) StreamProbeResults(probe *Probe, streamServer TcpingService_StreamProbeResultsServer) error {
	panic("not implemented")
}
