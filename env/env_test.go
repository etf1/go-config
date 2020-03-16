package env

import (
	"context"
	"os"
	"testing"

	"github.com/heetch/confita/backend"
	"github.com/stretchr/testify/require"
)

func TestEnvBackend(t *testing.T) {
	t.Run("NotFoundBecauseUnset", func(t *testing.T) {
		b := NewBackend()

		_, err := b.Get(context.Background(), "something that doesn't exist")
		require.Equal(t, backend.ErrNotFound, err)

	})

	t.Run("FoundEmpty", func(t *testing.T) {
		b := NewBackend()
		os.Setenv("TESTCONFIG", "")
		val, err := b.Get(context.Background(), "TESTCONFIG")
		require.NoError(t, err)
		require.Equal(t, "", string(val))
	})

	t.Run("ExactMatch", func(t *testing.T) {
		b := NewBackend()

		os.Setenv("TESTCONFIG1", "ok")
		val, err := b.Get(context.Background(), "TESTCONFIG1")
		require.NoError(t, err)
		require.Equal(t, "ok", string(val))
	})

	t.Run("DifferentCase", func(t *testing.T) {
		b := NewBackend()

		os.Setenv("TEST_CONFIG_2", "ok")
		val, err := b.Get(context.Background(), "test-config-2")
		require.NoError(t, err)
		require.Equal(t, "ok", string(val))
	})
}
