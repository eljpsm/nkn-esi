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
