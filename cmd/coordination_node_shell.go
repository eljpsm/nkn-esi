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
	// priceMap is the currently stored price map.
	priceMap = esi.PriceMap{}
	// resourceCharacteristics is the currently stored DER characteristics.
	resourceCharacteristics = esi.DerCharacteristics{}

	// receivedRegistrationForms are the currently stored registration forms.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)
	// registeredExchange is the public key of the engaged customer facility.
	//
	// As opposed to producers, there should only ever be one customer at any given time.
	registeredExchange = ""
	// registeredFacilities is a map of all other facilities registered in a producer role.
	registeredFacilities = make(map[string]bool)
	// priceMapOffers is a map of the current price map offers by uuid.
	priceMapOffers = make(map[string]*esi.PriceMapOffer)
	// producerPriceMaps are the price maps of the currently stored facilities engaged in a consumer role.
	producerPriceMaps = make(map[string]*esi.PriceMap)
	// producerCharacteristics are the characteristics of the currently stored facilities engaged in a consumer role.
	producerCharacteristics = make(map[string]*esi.DerCharacteristics)

	// auto accept details
	autoMoney = esi.Money{
		CurrencyCode: "NZD",
		Units:        100,
		Nanos:        0,
	}
	avoidMoney = esi.Money{
		CurrencyCode: "NZD",
		Units:        1000,
		Nanos:        0,
	}
	autoPrice = esi.PriceParameters{
		AlwaysBuyBelowPrice: &autoMoney,
		AvoidBuyOverPrice:   &avoidMoney,
	}
)

// coordinationNodeShell is the main shell of a coordination node.
func coordinationNodeShell() {
	logName := strings.TrimSuffix(coordinationNodePath, filepath.Ext(coordinationNodePath)) + logSuffix
	logFile, _ := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(logFile)
	log.SetLevel(log.InfoLevel)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go coordinationNodeMessageReceiver()
	wg.Add(2)
	go coordinationNodeInputReceiver()

	wg.Wait()
}
