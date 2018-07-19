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
	if schedule.Schedule == "" {
		return nil, E_EMPTY_SCHEDULE
	}

	if schedule.Id == "" {
		if schedule.Probe == nil || schedule.Probe.Id == "" {
			return nil, E_PROBE_ID
		}

		probe, probeExists := server.probes[schedule.Probe.Id]
		if probeExists {
			entryId, err := server.scheduler.AddFunc(schedule.Schedule, func() {
				result, err := probe.GetLatency()
				LOGGER.Info("Latency: %+v", result)
				server.resultQueue <- probeResult{result, err}
			})
			if err != nil {
				return nil, err
			}
			schedule.Id = string(entryId)
		} else {
			return nil, E_NO_PROBE
		}
	} else {
		_, probeExists := server.probes[schedule.Probe.Id]
		if schedule.Probe != nil && schedule.Probe.Id != "" && probeExists {
			// delete schedule then reschedule
		} else {

		}
	}
	return schedule, nil
}

func (server TcpingdServer) UnscheduleProbe(ctx context.Context, schedule *ProbeSchedule) (*Empty, error) {
	return nil, nil
}

func (server TcpingdServer) GetProbeResults(ctx context.Context, query *ProbeQuery) (*ProbeQueryResults, error) {
	return nil, nil
}

func (server TcpingdServer) StreamProbeResults(probe *Probe, streamServer TcpingService_StreamProbeResultsServer) error {
	return nil
}

func (server TcpingdServer) GetProbes(ctx context.Context, e *Empty) (*Probes, error) {
	var probes []*Probe

	server.probesLock.Lock()
	for k, v := range server.probes {
		probes = append(probes, &Probe{
			Id:    k,
			Label: "",
			Host:  v.DstIP,
			Port:  int32(v.DstPort),
		})
	}
	defer server.probesLock.Unlock()
	return &Probes{Probes: probes}, nil
}

func (server TcpingdServer) GetSchedules(ctx context.Context, e *Empty) (*ProbeSchedules, error) {
	return nil, nil
}
