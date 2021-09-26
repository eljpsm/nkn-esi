package cmd

import (
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// receivedRegistrationForms is a list of registration forms the user can then fill out.
var receivedRegistrationForms = []esi.DerFacilityRegistrationForm{}

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
