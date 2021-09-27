package cmd

import (
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// priceMapOffers is the currently stored price map offers.
	priceMapOffers = make(map[string]*esi.PriceMapOffer)

	// priceMapOfferFeedbacks is the currently stored price map offer feedbacks.
	priceMapOfferFeedbacks = make(map[string]*esi.PriceMapOfferFeedback)

	// priceMaps is the currently stored price map per negotiation.
	priceMaps = make(map[string]*esi.PriceMap)

	// receivedRegistrationForms is the currently stored registration forms.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)

	// producerFacilities is the currently stored facilities engaged as a producer role.
	producerFacilities = make(map[string]*esi.DerFacilityExchangeInfo)

	// consumerFacilities is the currently stored facilities engaged as a consumer role.
	consumerFacilities = make(map[string]*esi.DerFacilityExchangeInfo)
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
