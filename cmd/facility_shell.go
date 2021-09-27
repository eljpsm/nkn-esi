package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// facilityLoop is the main shell of a Facility.
func facilityShell() {
	logName := strings.TrimSuffix(facilityPath, filepath.Ext(facilityPath)) + logSuffix
	logFile, _ := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(logFile)
	log.SetLevel(log.InfoLevel)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go facilityMessageReceiver()
	wg.Add(2)
	go facilityInputReceiver()

	wg.Wait()
}
