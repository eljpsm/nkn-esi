/*
Copyright © 2021 Ecogy Energy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	// receivedRegistrationForms is a map of the currently stored registration forms.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)
	// registeredExchange is the public key of the engaged customer facility.
	//
	// As opposed to facilities, there should only ever be one customer at any given time.
	registeredExchange = ""
	// registeredFacilities is a map of all other facilities registered in a facility role.
	registeredFacilities = make(map[string]bool)
	// priceMapOffers is a map of the current price map offers by uuid.
	priceMapOffers = make(map[string]*esi.PriceMapOffer)
	// priceMapOfferStatus is a map of the status of stored price maps.
	priceMapOfferStatus = make(map[string]*esi.PriceMapOfferStatus)
	// facilityPriceMaps are the price maps of the currently stored facilities engaged in an exchange role.
	facilityPriceMaps = make(map[string]*esi.PriceMap)
	// facilityCharacteristics are the characteristics of the currently stored facilities engaged in a facility role.
	facilityCharacteristics = make(map[string]*esi.DerCharacteristics)

	// autoMoney is the money interface used for auto purchasing.
	autoMoney = esi.Money{
		CurrencyCode: "USD",
		Units:        100,
		Nanos:        0,
	}
	// avoidMoney is the money interface used for avoid purchasing. Currently, this is not used.
	avoidMoney = esi.Money{
		CurrencyCode: "USD",
		Units:        1000,
		Nanos:        0,
	}
	// autoPrice is the price parameters used for auto purchasing.
	autoPrice = esi.PriceParameters{
		AlwaysBuyBelowPrice: &autoMoney,
		AvoidBuyOverPrice:   &avoidMoney, // unused!
	}

	// voltageRange is the voltage range in volts.
	voltageRange = esi.SignedInt32Range{
		Min: 117,
		Max: 123,
	}
	// powerFactorRange is the power factor rage.
	powerFactorRange = esi.FloatRange{
		Min: 0.9,
		Max: 1.02,
	}
	// frequencyRange is the frequency range in hertz.
	frequencyRange = esi.SignedInt32Range{
		Min: 59,
		Max: 61,
	}
	//powerParameters is the expected power parameters.
	powerParameters = esi.PowerParameters{
		VoltageRange:     &voltageRange,
		PowerFactorRange: &powerFactorRange,
		FrequencyRange:   &frequencyRange,
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
	go coordinationNodeMessageReceiver() // receive incoming messages
	wg.Add(2)
	go coordinationNodeInputReceiver() // receive user input
	wg.Add(3)
	go coordinationNodePeriodicMessenger() // send regular information to any facilities

	wg.Wait()
}
