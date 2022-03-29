package writersreg

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertMutexesEqualAddress(t *testing.T, mutexA, mutexB *sync.Mutex) {
	t.Helper()

	if mutexA == nil && mutexB == nil {
		return
	}

	addressA := fmt.Sprintf("%p", mutexA)
	addressB := fmt.Sprintf("%p", mutexB)
	assert.Equal(t, addressA, addressB)
}
