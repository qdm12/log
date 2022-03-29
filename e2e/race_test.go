package log

import (
	"bytes"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/qdm12/log"
)

func Test_Logger_Writer_Race(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	parent := log.New(log.SetWriters(buffer))

	childA := parent.New()
	childB := parent.New()

	childAA := childA.New()
	childBA := childB.New()

	loggers := []*log.Logger{
		parent,
		childA,
		childB,
		childAA,
		childBA,
	}

	readyWait := new(sync.WaitGroup)
	readyWait.Add(len(loggers))

	doneWait := new(sync.WaitGroup)
	doneWait.Add(len(loggers))

	// run for 50ms
	stopCh := make(chan struct{})
	go func() {
		const timeout = 50 * time.Millisecond
		readyWait.Wait()
		<-time.After(timeout)
		close(stopCh)
	}()

	for _, logger := range loggers {
		go func(logger *log.Logger) {
			defer doneWait.Done()
			readyWait.Done()
			readyWait.Wait()

			for {
				select {
				case <-stopCh:
					return
				default:
				}

				// test relies on the -race detector
				// to detect concurrent writes to the buffer.
				logger.Info("x")
			}
		}(logger)
	}
	doneWait.Wait()
}

func Test_Logger_Settings_Race(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	logger := log.New(log.SetWriters(buffer))

	workers := runtime.NumCPU()

	readyWait := new(sync.WaitGroup)
	readyWait.Add(workers)

	doneWait := new(sync.WaitGroup)
	doneWait.Add(workers)

	// run for 50ms
	stopCh := make(chan struct{})
	go func() {
		const timeout = 50 * time.Millisecond
		readyWait.Wait()
		<-time.After(timeout)
		close(stopCh)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			defer doneWait.Done()
			readyWait.Done()
			readyWait.Wait()

			for {
				select {
				case <-stopCh:
					return
				default:
				}

				// test relies on the -race detector
				// to detect concurrent writes to the buffer.
				logger.Info("x")
				logger.Patch(log.SetLevel(log.LevelInfo))
			}
		}()
	}
	doneWait.Wait()
}
