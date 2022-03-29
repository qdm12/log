package writersreg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Registry_RegisterWriters(t *testing.T) {
	t.Parallel()

	writers := []io.Writer{os.Stderr, os.Stdout}

	registry := NewRegistry()

	mutexesA := registry.RegisterWriters(writers)
	require.Equal(t, len(writers), len(mutexesA))
	for i, writer := range writers {
		writerAddress := fmt.Sprintf("%p", writer)
		registryMutex := registry.writerAddressToMutex[writerAddress]
		assertMutexesEqualAddress(t, mutexesA[i], registryMutex)
	}

	mutexesB := registry.RegisterWriters(writers)
	require.Equal(t, len(mutexesA), len(mutexesB))
	for i, mutexB := range mutexesB {
		mutexA := mutexesA[i]
		assertMutexesEqualAddress(t, mutexA, mutexB)
	}

	writers = []io.Writer{bytes.NewBuffer(nil), bytes.NewBuffer(nil)}
	mutexesC := registry.RegisterWriters(writers)
	require.Equal(t, len(writers), len(mutexesC))
	for i, writer := range writers {
		writerAddress := fmt.Sprintf("%p", writer)
		registryMutex := registry.writerAddressToMutex[writerAddress]
		assertMutexesEqualAddress(t, mutexesC[i], registryMutex)
	}

	assert.Len(t, registry.writerAddressToMutex, 4)
}

func Test_Registry_registerWriter(t *testing.T) {
	t.Parallel()

	t.Run("nil writer", func(t *testing.T) {
		t.Parallel()

		assert.PanicsWithValue(t, "writer cannot be nil", func() {
			registry := &Registry{}
			registry.registerWriter(nil)
		})
	})

	t.Run("io.Discard writer", func(t *testing.T) {
		t.Parallel()

		registry := &Registry{}
		mutex := registry.registerWriter(io.Discard)
		assert.Nil(t, mutex)
	})

	t.Run("writer already registered", func(t *testing.T) {
		t.Parallel()

		writer := bytes.NewBuffer(nil)
		writerAddress := fmt.Sprintf("%p", writer)

		existingMutex := new(sync.Mutex)

		registry := &Registry{
			writerAddressToMutex: map[string]*sync.Mutex{
				writerAddress: existingMutex,
			},
		}

		mutex := registry.registerWriter(writer)

		assertMutexesEqualAddress(t, existingMutex, mutex)

		expectedRegistry := &Registry{
			writerAddressToMutex: map[string]*sync.Mutex{
				writerAddress: existingMutex,
			},
		}
		assert.Equal(t, expectedRegistry, registry)
	})

	t.Run("writer not registered", func(t *testing.T) {
		t.Parallel()

		writer := bytes.NewBuffer(nil)
		writerAddress := fmt.Sprintf("%p", writer)

		registry := &Registry{
			writerAddressToMutex: map[string]*sync.Mutex{},
		}

		mutex := registry.registerWriter(writer)

		registryMutex := registry.writerAddressToMutex[writerAddress]
		assertMutexesEqualAddress(t, registryMutex, mutex)

		expectedRegistry := &Registry{
			writerAddressToMutex: map[string]*sync.Mutex{
				writerAddress: registryMutex,
			},
		}
		assert.Equal(t, expectedRegistry, registry)
	})
}
