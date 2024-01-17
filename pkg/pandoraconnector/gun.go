package pandoraconnector

import (
	"fmt"

	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/register"

	"github.com/lueurxax/e2e/common"
	"github.com/lueurxax/e2e/pkg/models"
)

// S3Gun is gun name for test s3
const S3Gun = "s3_gun"

// gunConfig config of s3 gun
type gunConfig struct {
	tester models.Tester
	client string
}

// Gun is S3Gun structure
type Gun struct {
	// Configured on Bind, before shooting
	aggr core.Aggregator // May be your custom Aggregator.
	core.GunDeps
	configurator gunConfigurator
}

// Bind gun to engine
func (g *Gun) Bind(aggr core.Aggregator, deps core.GunDeps) (err error) {
	g.aggr = aggr
	g.GunDeps = deps
	return nil
}

// Shoot ammo with gun
func (g *Gun) Shoot(ammo core.Ammo) {
	customAmmo, ok := ammo.(*Ammo)
	if !ok {
		g.Log.Error("invalid structure of ammo")
	}
	g.shoot(customAmmo)
}

func (g *Gun) shoot(ammo *Ammo) {
	code := 0
	conf := g.configurator.GetGunConfig()
	sample := common.NewRequestData(conf.tester.MethodName())
	defer func() {
		sample.SetProtoCode(code)
		g.aggr.Report(sample)
	}()

	_, err := conf.tester.Run(g.Ctx, conf.client, ammo.Params, &models.Options{
		Conf: ammo.Conf,
	})

	if err != nil {
		fmt.Printf("FATAL: %s", err)
		code = 500
		return
	}
	code = 200
}

// RegisterGun construct new gun
func RegisterGun(configurator gunConfigurator) {
	register.Gun(S3Gun, func() core.Gun {
		return &Gun{configurator: configurator}
	})
}

type gunConfManager struct {
	config gunConfig
}

func (g *gunConfManager) SetClient(client string) {
	g.config.client = client
}

func (g *gunConfManager) SetTester(tester models.Tester) {
	g.config.tester = tester
}

func (g *gunConfManager) GetGunConfig() gunConfig {
	return g.config
}

func newGunConf() gunConfigurator {
	return &gunConfManager{config: gunConfig{}}
}
