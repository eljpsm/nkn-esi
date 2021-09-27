package cmd

import (
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// priceMapOffers are the currently stored price map offers.
	priceMapOffers = make(map[uuid.UUID]*esi.PriceMapOffer)

	// priceMapOfferFeedbacks are the currently stored price map offer feedbacks.
	priceMapOfferFeedbacks = make(map[uuid.UUID]*esi.PriceMapOfferFeedback)

	// priceMap is the currently stored price map.
	priceMap = esi.PriceMap{}
	// isPriceMapAccepted is a simple flag showing whether the price map has been accepted.
	isPriceMapAccepted = false

	// resourceCharacteristics is the currently stored DER characteristics.
	resourceCharacteristics = esi.DerCharacteristics{}

	// receivedRegistrationForms are the currently stored registration forms.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)

	// customerFacilities is a map of all other facilities registered in a customer role.
	customerFacilities = make(map[string]bool)
	// customerPriceMapOffers is a map of the current price map offers by public key.
	customerPriceMapOffers = make(map[uuid.UUID]*esi.PriceMapOffer)
	// producerFacilities is a map of all other facilities registered in a producer role.
	producerFacilities = make(map[string]bool)
	// producerPriceMaps are the price maps of the currently stored facilities engaged in a consumer role.
	producerPriceMaps = make(map[string]*esi.PriceMap)
	// producerCharacteristics are the characteristics of the currently stored facilities engaged in a consumer role.
	producerCharacteristics = make(map[string]*esi.DerCharacteristics)

	// auto accept details
	autoMoney = esi.Money{
		CurrencyCode: "NZD",
		Units: 100,
		Nanos: 0,
	}
	avoidMoney = esi.Money{
		CurrencyCode: "NZD",
		Units: 1000,
		Nanos: 0,
	}
	autoPrice = esi.PriceParameters{
		AlwaysBuyBelowPrice: &autoMoney,
		AvoidBuyOverPrice: &avoidMoney,
	}
)

const (
	uuidHigh = 127
	uuidLow  = 0
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
