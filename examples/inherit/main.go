package main

import (
	"github.com/qdm12/log"
)

func main() {
	loggerA := log.New(log.SetComponent("A"))
	loggerB := loggerA.New(log.SetComponent("B"))
	loggerA.Info("my message")
	// 2022-03-29T07:35:08Z INFO [A] my message
	loggerB.Info("my message")
	// 2022-03-29T07:35:08Z INFO [B] my message
}
