package testerspool

import (
	"github.com/lueurxax/e2e/common"
	"github.com/lueurxax/e2e/pkg/models"
)

// TestersPool readonly pool of tester interfaces
type TestersPool interface {
	Get(key string) (tester models.Tester, err error)
}

type pool struct {
	storage map[string]models.Tester
}

// Get tester from pool
func (t *pool) Get(key string) (tester models.Tester, err error) {
	var ok bool
	tester, ok = t.storage[key]
	if !ok {
		err = common.ErrUnknownMethod(key)
	}
	return
}

// NewTestersPool constructor for testers pool
func NewTestersPool(testers []models.Tester) TestersPool {
	p := map[string]models.Tester{}
	for _, tester := range testers {
		p[tester.MethodName()] = tester
	}
	return &pool{storage: p}
}
