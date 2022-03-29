package writersreg

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRegistry(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	expectedRegistry := &Registry{
		writerAddressToMutex: make(map[string]*sync.Mutex, 1),
	}

	assert.Equal(t, expectedRegistry, registry)
}
