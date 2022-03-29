package main

import (
	"os"
	"time"

	"github.com/qdm12/log"
)

func main() {
	logger := log.New(
		log.SetLevel(log.LevelDebug),
		log.SetTimeFormat(time.RFC822),
		log.SetWriters(os.Stdout, os.Stderr),
		log.SetComponent("module"),
		log.SetCallerFile(true),
		log.SetCallerFunc(true),
		log.SetCallerLine(true))
	logger.Info("my message")
}
