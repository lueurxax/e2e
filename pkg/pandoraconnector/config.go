package pandoraconnector

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
	pandoraconfig "github.com/yandex/pandora/core/config"
	"github.com/yandex/pandora/core/engine"

	"github.com/lueurxax/e2e/common"
)

// TODO will replaced with constructor for engine.Config
func initConfig(conf common.StressLoad) (engineConf *engine.Config, err error) {
	duration := strconv.Itoa(conf.Duration) + "s"

	id := uuid.NewV4()

	confMap := &map[string][]map[string]interface{}{
		"pools": {
			{
				"id": id.String(),
				"gun": map[string]interface{}{
					"type": S3Gun,
				},
				"ammo": map[string]interface {
				}{
					"type": S3Provider,
				},
				"result": map[string]interface {
				}{
					"type": s3Aggregator,
				},
				"rps": map[string]interface{}{
					"duration": duration,
					"type":     "line",
					"from":     conf.From,
					"to":       conf.To,
				},
				"startup": map[string]interface {
				}{
					"type":  "once",
					"times": conf.Instances,
				},
			},
		},
	}

	confStruct := &struct {
		Engine engine.Config `config:",squash"`
	}{}

	if err := pandoraconfig.DecodeAndValidate(confMap, confStruct); err != nil {
		return nil, err
	}
	return &confStruct.Engine, nil
}
