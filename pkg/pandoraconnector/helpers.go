package pandoraconnector

import (
	"context"
	"fmt"
	"time"

	"github.com/yandex/pandora/core/engine"
	"github.com/yandex/pandora/lib/monitoring"
	"go.uber.org/zap"
)

func runEngine(ctx context.Context, engine *engine.Engine, errs chan error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errs <- engine.Run(ctx)
}

func newEngineMetrics() engine.Metrics {
	return engine.Metrics{
		Request:        monitoring.NewCounter("engine_Requests"),
		Response:       monitoring.NewCounter("engine_Responses"),
		InstanceStart:  monitoring.NewCounter("engine_UsersStarted"),
		InstanceFinish: monitoring.NewCounter("engine_UsersFinished"),
	}
}

func startReport(m engine.Metrics) (cancel chan struct{}) {
	requests := m.Request.Get()
	responses := m.Response.Get()
	cancel = make(chan struct{})
	go func(cancel chan struct{}) {
		var requestsNew, responsesNew int64
		// TODO(skipor): there is no guarantee, that we will run exactly after 1 second.
		// So, when we get 1 sec +-10ms, we getting 990-1010 calculate intervals and +-2% RPS in reports.
		// Consider using rcrowley/go-metrics.Meter.
	exit:
		for {
			select {
			case <-time.NewTicker(1 * time.Second).C:
				requestsNew = m.Request.Get()
				responsesNew = m.Response.Get()
				rps := responsesNew - responses
				reqps := requestsNew - requests
				activeUsers := m.InstanceStart.Get() - m.InstanceFinish.Get()
				activeRequests := requestsNew - responsesNew
				fmt.Printf(
					"[ENGINE] %d resp/s; %d req/s; %d users; %d active\n",
					rps, reqps, activeUsers, activeRequests)

				requests = requestsNew
				responses = responsesNew
			case <-cancel:
				break exit
			}
		}
	}(cancel)
	return
}

func newLogger() *zap.Logger {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.Level.SetLevel(zap.DebugLevel)
	zapLog, err := zapConf.Build(zap.AddCaller())
	if err != nil {
		zap.L().Fatal("Logger build failed", zap.Error(err))
	}
	return zapLog
}
