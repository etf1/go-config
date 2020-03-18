package prefix

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BackendMock struct{}

func NewBackendMock() *BackendMock {
	return &BackendMock{}
}

func (b *BackendMock) Get(ctx context.Context, key string) ([]byte, error) {
	return []byte(key), nil
}

func (b *BackendMock) Name() string {
	return "test_backend"
}

func TestNewBackend(t *testing.T) {
	// Given
	prefix := "my_prefix"
	backendMock := NewBackendMock()

	// When
	prefixBackend := NewBackend(prefix, backendMock)

	// Then
	assert := assert.New(t)

	assert.Equal("prefix_test_backend", prefixBackend.Name())

	assert.Equal(prefix, prefixBackend.prefix)
	assert.Equal(backendMock, prefixBackend.backend)
}

func TestGetWhenPrefixed(t *testing.T) {
	// Given
	ctx := context.Background()

	prefix := "my_prefix"
	backendMock := NewBackendMock()

	prefixBackend := NewBackend(prefix, backendMock)

	// When
	result, err := prefixBackend.Get(ctx, "my_key")

	// Then
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal("my_prefix_my_key", string(result))
}

func TestGetWhenPrefixedWithAnotherDelimiter(t *testing.T) {
	// Given
	ctx := context.Background()

	prefix := "my_prefix"
	backendMock := NewBackendMock()

	prefixBackend := NewBackend(prefix, backendMock, WithDelimiter("."))

	// When
	result, err := prefixBackend.Get(ctx, "my_key")

	// Then
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal("my_prefix.my_key", string(result))
}

func TestName(t *testing.T) {
	// Given
	prefix := "my_prefix"
	backendMock := NewBackendMock()

	prefixBackend := NewBackend(prefix, backendMock, WithDelimiter("."))

	// When - Then
	assert.Equal(t, "prefix_test_backend", prefixBackend.Name())
}
