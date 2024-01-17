package pandoraconnector

import (
	"context"

	"github.com/yandex/pandora/core"
	pandoraprov "github.com/yandex/pandora/core/provider"
	"github.com/yandex/pandora/core/register"
	"go.uber.org/zap"

	"github.com/lueurxax/e2e/pkg/models"
)

// S3Provider if provider name for stress test s3
const S3Provider = "rest_provider"

// Provider is routine that generates ammo for Instance shoots.
// A Provider MUST be goroutine safe.
type Provider interface {
	// Run starts provider routine of ammo  generation.
	// Blocks until ammo finish, error or context cancel.
	// Run MUST be called only once. Run SHOULD be called before Acquire or Release calls, but
	// MAY NOT because of goroutine races.
	// In case of ctx cancel, SHOULD return nil, but MAY ctx.Err(), or error caused ctx.Err()
	// in terms of github.com/pkg/errors.Cause.
	Run(ctx context.Context, deps core.ProviderDeps) error
	// Acquire acquires ammo for shoot. Acquire SHOULD be lightweight, so Instance can Shoot as
	// soon as possible. That means ammo format parsing SHOULD be done in Provider Run goroutine,
	// but acquire just takes ammo from ready queue.
	// Ok false means that shooting MUST be stopped because ammo finished or shooting is canceled.
	// Acquire MAY be called before Run, but SHOULD block until Run is called.
	Acquire() (ammo core.Ammo, ok bool)
	// Release notifies that ammo usage is finished, and it can be reused.
	// Instance MUST NOT retain references to released ammo.
	Release(ammo core.Ammo)
}

type provider struct {
	pandoraprov.AmmoQueue
	*core.ProviderDeps
	providerConfigurator
}

// Run starts provider routine of ammo  generation.
// Blocks until ammo finish, error or context cancel.
// Run MUST be called only once. Run SHOULD be called before Acquire or Release calls, but
// MAY NOT because of goroutine races.
// In case of ctx cancel, SHOULD return nil, but MAY ctx.Err(), or error caused ctx.Err()
// in terms of github.com/pkg/errors.Cause.
func (p *provider) Run(ctx context.Context, deps core.ProviderDeps) (err error) {
	p.ProviderDeps = &deps
	p.Log.Info("run provider")
	defer close(p.OutQueue)
	for i, param := range p.providerConfigurator.GetParams() {
		select {
		case p.OutQueue <- &Ammo{
			Params: param,
			Conf:   p.providerConfigurator.Conf(),
		}:
		case <-ctx.Done():
			p.Log.Debug("Provider run context is Done", zap.Int("decoded", i+1))
			return nil
		}
	}
	<-ctx.Done()
	return
}

// RegisterProvider register new s3 provider
func RegisterProvider(
	providerConfigurator providerConfigurator,
) {
	register.Provider(S3Provider, func() core.Provider {
		newAmmo := func() core.Ammo { return map[string]interface{}{} }
		p := &provider{
			AmmoQueue:            *pandoraprov.NewAmmoQueue(newAmmo, pandoraprov.DefaultAmmoQueueConfig()),
			providerConfigurator: providerConfigurator,
		}
		return p
	})
}

type provConfigManager struct {
	params []models.StateSelector
	opts   *models.Options
}

func newProvConfig() providerConfigurator {
	return &provConfigManager{}
}

func (p *provConfigManager) SetParameters(params []models.StateSelector, opts *models.Options) {
	p.params = params
	p.opts = opts
}

func (p *provConfigManager) GetParams() (params []models.StateSelector) {
	return p.params
}

func (p *provConfigManager) Conf() (conf *models.Test) {
	return p.opts.Conf
}
