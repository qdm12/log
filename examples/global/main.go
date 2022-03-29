package main

import (
	"bytes"
	"time"

	"github.com/qdm12/log"
)

func main() {
	writer := bytes.NewBuffer(nil)

	timer := time.NewTimer(time.Second)

	const parallelism = 2
	for i := 0; i < parallelism; i++ {
		go func() {
			logger := log.New(log.SetWriters(writer))
			for {
				select {
				case <-timer.C:
					return
				default:
					logger.Info("my message")
				}
			}
		}()
	}
}
