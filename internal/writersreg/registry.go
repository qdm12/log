package writersreg

import (
	"sync"
)

func NewRegistry() *Registry {
	const initialWriterCapacity = 1
	return &Registry{
		writerAddressToMutex: make(map[string]*sync.Mutex, initialWriterCapacity),
	}
}

type Registry struct {
	mutex                sync.RWMutex
	writerAddressToMutex map[string]*sync.Mutex
}
