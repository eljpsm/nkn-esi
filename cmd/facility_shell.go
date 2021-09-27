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
	// priceMapOffers are the currently stored price map offers.
	priceMapOffers = make(map[string]*esi.PriceMapOffer)

	// priceMapOfferFeedbacks are the currently stored price map offer feedbacks.
	priceMapOfferFeedbacks = make(map[string]*esi.PriceMapOfferFeedback)

	// priceMap is the currently stored price map.
	priceMap = esi.PriceMap{}

	// receivedRegistrationForms are the currently stored registration forms.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)

	// consumerFacilitiesPriceMaps are the price maps of the currently stored facilities engaged in a consumer role.
	consumerFacilitiesPriceMaps = make(map[string]*esi.PriceMap)
	// consumerFacilitiesCharacteristics are the characteristics of the currently stored facilities engaged in a consumer role.
	consumerFacilitiesCharacteristics = make(map[string]*esi.DerCharacteristics)
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
