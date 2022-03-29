package main

import "github.com/qdm12/log"

func main() {
	logger := log.New()
	logger.Warnf("message number %d", 1)
	// 2022-03-29T07:40:12Z WARN message number 1
}
