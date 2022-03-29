package writersreg

import (
	"fmt"
	"io"
	"sync"
)

func (r *Registry) RegisterWriters(writers []io.Writer) (mutexes []*sync.Mutex) {
	mutexes = make([]*sync.Mutex, len(writers))

	for i, writer := range writers {
		mutexes[i] = r.registerWriter(writer)
	}

	return mutexes
}

func (r *Registry) registerWriter(writer io.Writer) (mutex *sync.Mutex) {
	if writer == nil {
		panic("writer cannot be nil")
	}

	if writer == io.Discard {
		return nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	writerAddress := fmt.Sprintf("%p", writer)

	mutex, ok := r.writerAddressToMutex[writerAddress]
	if ok {
		// writer already registered
		return mutex
	}

	mutex = new(sync.Mutex)
	r.writerAddressToMutex[writerAddress] = mutex
	return mutex
}
