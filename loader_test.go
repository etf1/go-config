package config_test

import (
	"context"
	"testing"

	"github.com/heetch/confita/backend"
	"github.com/stretchr/testify/assert"

	"github.com/etf1/go-config"
)

type Conf struct {
	Param1 string `config:"param1"`
	Param2 string `config:"param2"`
	Param3 string `config:"param3"`
}

func TestLoad(t *testing.T) {
	baseConfig := &Conf{
		Param1: "param1 default value",
		Param3: "param3 default value",
	}
	backends := []backend.Backend{
		mapBackend("source1", map[string]string{"param1": "111", "param2": "222", "param3": "333"}),
		mapBackend("source1", map[string]string{"param1": "AAA", "param2": "BBB"}),
		mapBackend("source3", map[string]string{"param3": ""}),
		mapBackend("source4", map[string]string{"param999": "param999 value"}),
	}
	expectedConfig := &Conf{
		Param1: "AAA",
		Param2: "BBB",
		Param3: "",
	}

	assertConfig(t, expectedConfig, baseConfig, backends)
}

func assertConfig(t *testing.T, expectedConfig *Conf, baseConfig *Conf, backends []backend.Backend) {
	l := config.NewConfigLoader(backends...)
	assert.Nil(t, l.Load(context.Background(), baseConfig))
	assert.Equal(t,
		expectedConfig,
		baseConfig,
	)
}

func mapBackend(name string, values map[string]string) backend.Backend {
	return backend.Func(name, func(ctx context.Context, key string) ([]byte, error) {
		if v, ok := values[key]; ok {
			return []byte(v), nil
		}
		return nil, backend.ErrNotFound
	})
}
