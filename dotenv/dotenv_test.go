package dotenv_test

import (
	"context"
	"os"
	"testing"

	"github.com/etf1/go-config"
	"github.com/stretchr/testify/assert"

	"github.com/etf1/go-config/dotenv"
)

func TestNewBackend_UnknowFile(t *testing.T) {
	assert.Panics(t, func() {
		dotenv.NewBackend("unknow_file")
	})
}

func TestNewBackend_Panic(t *testing.T) {
	t.Run("Is a directory", func(t *testing.T) {
		assert.PanicsWithError(t, "read ./testdata: is a directory", func() {
			dotenv.NewBackend("./testdata")
		})
	})

	t.Run("Permission denied", func(t *testing.T) {
		_ = os.Chmod("./testdata/.env.denied", 000)
		defer os.Chmod("./testdata/.env.denied", 644)
		assert.PanicsWithError(t, "open ./testdata/.env.denied: permission denied", func() {
			dotenv.NewBackend("./testdata/.env.denied")
		})
	})
}

func TestNewBackend(t *testing.T) {
	b := dotenv.NewBackend("testdata/.env")

	t.Run("Key format unexpected", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "invalid key with space")
		assert.EqualError(t, err, "dotenv variable format expected \\w+, \"invalid key with space\" given")
		assert.Equal(t, "", string(result))
	})

	t.Run("Key Not Found", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "test")
		assert.EqualError(t, err, "key not found")
		assert.Equal(t, "", string(result))
	})

	t.Run("Key Matched", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "my_string")
		assert.NoError(t, err)
		assert.Equal(t, "my value", string(result))
	})

	t.Run("Value Multiline Matched", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "multi")
		assert.NoError(t, err)

		expectedValue := `a
b`
		assert.Equal(t, expectedValue, string(result))
	})
}

func TestNewBackend_Multiple(t *testing.T) {
	b := dotenv.NewBackend("testdata/.env", "testdata/.env.other")

	t.Run("Key format unexpected", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "invalid key with space")
		assert.EqualError(t, err, "dotenv variable format expected \\w+, \"invalid key with space\" given")
		assert.Equal(t, "", string(result))
	})

	t.Run("Key Not Found", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "test")
		assert.EqualError(t, err, "key not found")
		assert.Equal(t, "", string(result))
	})

	t.Run("Key Matched", func(t *testing.T) {
		result, err := b.Get(context.TODO(), "my_string")
		assert.NoError(t, err)
		assert.Equal(t, "my other value", string(result))
	})
}

func TestGetBackend_Panic(t *testing.T) {
	t.Run("File denied", func(t *testing.T) {
		_ = os.Chmod("./testdata/.env.denied", 000)
		defer os.Chmod("./testdata/.env.denied", 644)

		assert.PanicsWithError(t, "open testdata/.env.denied: permission denied", func() {
			dotenv.GetBackends("testdata/.env.denied")
		})
	})
}

func TestGetBackend(t *testing.T) {
	t.Run("One Backend", func(t *testing.T) {
		backends := dotenv.GetBackends("testdata/.env")
		assert.Len(t, backends, 1)
	})

	t.Run("Two Backend", func(t *testing.T) {
		backends := dotenv.GetBackends("testdata/.env", "testdata/.env.other")
		assert.Len(t, backends, 2)
	})

	t.Run("File Not found", func(t *testing.T) {
		backends := dotenv.GetBackends("testdata/.env", "unexisting.file")
		assert.Len(t, backends, 1)
	})

	t.Run("Directory", func(t *testing.T) {
		backends := dotenv.GetBackends("testdata/.env", "./testdata")
		assert.Len(t, backends, 1)
	})
}

// This function will replace os.Args in order to simulate only wanted flags
func replaceOSArgs(flags ...string) func() {
	originalArgs := os.Args
	os.Args = append([]string{originalArgs[0]}, flags...)
	return func() {
		os.Args = originalArgs
	}
}

func TestGetBackendsFromFlag(t *testing.T) {
	defer replaceOSArgs("-config-env-files=testdata/.env,testdata/.env.other")()
	backends := dotenv.GetBackendsFromFlag()
	assert.Len(t, backends, 2)
}

type Conf struct {
	Param1 string `config:"my_string"`
	Param2 bool   `config:"my_bool"`
	Param3 string `config:"param3"`
}

func TestNewBackend_MultipleFiles(t *testing.T) {
	defer replaceOSArgs("-config-env-files=testdata/.env,testdata/.env.other")()
	baseConfig := &Conf{
		Param1: "param1 default value",
		Param3: "param3 default value",
	}
	l := config.NewDefaultConfigLoader()
	assert.Nil(t, l.Load(context.Background(), baseConfig))
	assert.Equal(t, "my other value", baseConfig.Param1)
	assert.True(t, baseConfig.Param2)
	assert.Equal(t, "param3 default value", baseConfig.Param3)
}
