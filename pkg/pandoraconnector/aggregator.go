package pandoraconnector

import (
	"context"

	"github.com/yandex/pandora/core"

	"git.proksy.io/golang/e2e/common"
	"git.proksy.io/golang/e2e/pkg/log"
)

const s3Aggregator = "rest_aggregator"

// Aggregator is routine that aggregates Samples from all Pool Instances.
// Usually aggregator is shooting result reporter, that writes Reported Samples
// to DataSink in machine readable format for future analysis.
// An Aggregator MUST be goroutine safe.
// GunDeps are passed to Gun before Instance Run.
type Aggregator interface {
	// Run starts aggregator routine of handling Samples. Blocks until fail or context cancel.
	// Run MUST be called only once. Run SHOULD be called before Report calls, but MAY NOT because
	// of goroutine races.
	// In case of ctx cancel, SHOULD return nil, but MAY ctx.Err(), or error caused ctx.Err()
	// in terms of github.com/pkg/errors.Cause.
	// In case of any dropped RequestData (unhandled because of RequestData queue overflow) Run SHOULD return
	// error describing how many samples were dropped.
	Run(ctx context.Context, deps core.AggregatorDeps) error
	// Report reports sample to aggregator. SHOULD be lightweight and not blocking,
	// so Instance can Shoot as soon as possible.
	// That means, that RequestData encode and reporting SHOULD NOT be done in caller goroutine,
	// but SHOULD in Aggregator Run goroutine.
	// If Aggregator can't handle Reported RequestData without blocking, it SHOULD just drop it.
	// Reported Samples MAY just be dropped, after context cancel.
	// Reported RequestData MAY be reused for efficiency, so caller MUST NOT retain reference to RequestData.
	// Report MAY be called before Aggregator Run. Report MAY be called after Run finish, in case of
	// Pool Run cancel.
	// Aggregator SHOULD Return RequestData if it implements BorrowedSample.
	Report(s core.Sample)
}

type aggregator struct {
	log   log.Logger
	sink  chan core.Sample
	meter common.Meter
}

func (s *aggregator) Run(ctx context.Context, deps core.AggregatorDeps) error {
loop:
	for {
		select {
		case sample := <-s.sink:
			s.handle(sample)
		case <-ctx.Done():
			break loop
		}
	}
	for {
		// Context is done, but we should read all data from sink.
		select {
		case r := <-s.sink:
			s.handle(r)
		default:
			return nil
		}
	}
}

func (s *aggregator) Report(sample core.Sample) {
	s.sink <- sample
}

func (s *aggregator) handle(sample core.Sample) {
	data := sample.(*common.RequestData)
	s.meter.AddRequest(data)
}

// News3Aggregator construct new S3 Aggregator
func News3Aggregator(log log.Logger, metrics common.Meter) func() Aggregator {
	return func() Aggregator {
		return &aggregator{log: log, sink: make(chan core.Sample, 128), meter: metrics}
	}
}
