package pandoraconnector

import (
	"context"
	"fmt"
	"time"

	"github.com/yandex/pandora/core/engine"
	"github.com/yandex/pandora/core/plugin/pluginconfig"
	"github.com/yandex/pandora/core/register"
	"github.com/yandex/pandora/core/schedule"
	"go.uber.org/zap"

	"git.proksy.io/golang/e2e/common"
	"git.proksy.io/golang/e2e/pkg/log"
	"git.proksy.io/golang/e2e/pkg/models"
)

const compositeScheduleKey = "composite"

// PandoraConnector for connect to pandora
type PandoraConnector interface {
	Register(metrics common.Meter)
	Start(
		ctx context.Context,
		client string,
		tester models.Tester,
		params []models.StateSelector,
		opts *models.Options,
	) (newParams []models.StateSelector, err error)
}

type gunConfigurator interface {
	SetTester(tester models.Tester)
	SetClient(client string)
	GetGunConfig() gunConfig
}

type providerConfigurator interface {
	SetParameters(params []models.StateSelector, opts *models.Options)
	GetParams() (params []models.StateSelector)
	Conf() (conf *models.Test)
}

type connector struct {
	gunConfigurator
	providerConfigurator
	logger        log.Logger
	engineMetrics engine.Metrics
}

func (c *connector) Register(metrics common.Meter) {
	register.Aggregator(s3Aggregator, News3Aggregator(c.logger, metrics))

	register.Limiter("line", schedule.NewLineConf)
	register.Limiter("const", schedule.NewConstConf)
	register.Limiter("once", schedule.NewOnceConf)
	register.Limiter("unlimited", schedule.NewUnlimitedConf)
	register.Limiter(compositeScheduleKey, schedule.NewCompositeConf)

	// Required for decoding plugins. Need to be added after Composite Schedule hacky hook.
	pluginconfig.AddHooks()

	// Custom imports. Integrate your custom types into configuration system.
	RegisterProvider(c.providerConfigurator)

	// Register gun
	RegisterGun(c.gunConfigurator)
}

// Start stress test with pandora
func (c *connector) Start(
	ctx context.Context,
	client string,
	tester models.Tester,
	params []models.StateSelector,
	opts *models.Options,
) (newParams []models.StateSelector, err error) {
	c.providerConfigurator.SetParameters(params, opts)
	c.gunConfigurator.SetClient(client)
	c.gunConfigurator.SetTester(tester)

	cancelReport := startReport(c.engineMetrics)

	zapLogger := newLogger()
	zap.ReplaceGlobals(zapLogger)
	zap.RedirectStdLog(zapLogger)

	var conf *engine.Config
	conf, err = initConfig(*opts.Conf.StressLoad)
	if err != nil {
		return
	}

	pandora := engine.New(zapLogger, c.engineMetrics, *conf)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error)
	go runEngine(ctx, pandora, errs)

	// waiting for signal or error message from engine
exit:
	for err = range errs {
		switch err {
		case nil:
			c.logger.Info("Pandora engine successfully finished it's work")
			break exit
		case err:
			const awaitTimeout = 3 * time.Second
			c.logger.Info(
				fmt.Sprintf(
					"engine run failed. Awaiting started tasks. error: %s, timeout: %d",
					err, awaitTimeout),
			)
			cancel()
			time.AfterFunc(awaitTimeout, func() {
				c.logger.Error(fmt.Errorf("engine tasks timeout exceeded"))
			})
			pandora.Wait()
			c.logger.Error(fmt.Errorf("engine run failed. Pandora graceful shutdown successfully finished"))
			return
		}
	}
	c.logger.Info("Engine run successfully finished")
	close(cancelReport)
	return params, nil
}

// NewConnector construct and register pandora instance
func NewConnector(logger log.Logger) PandoraConnector {
	// CreateLottery engine metrics
	m := newEngineMetrics()
	return &connector{
		logger:               logger,
		providerConfigurator: newProvConfig(),
		gunConfigurator:      newGunConf(),
		engineMetrics:        m,
	}
}
