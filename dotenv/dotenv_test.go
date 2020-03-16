package dotenv_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/etf1/config/dotenv"
)

func TestNewBackend_UnknowFile(t *testing.T) {
	assert.Panics(t, func() {
		dotenv.NewBackend("unknow_file")
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
}
