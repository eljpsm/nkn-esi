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
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// coordinationNodePeriodicMessenger sends information at a regular period.
func coordinationNodePeriodicMessenger() {
	// Expected price low of any random price.
	priceLow := 20
	// Expected price high of any random price.
	priceHigh := 35

	for {
		if len(registeredFacilities) > 0 {

			// Explicit look at just facilities.
			for publicKey := range registeredFacilities {
				// Send price datum.
				//
				// At the moment, all prices are random within some range. However, in practice, this would probably
				// take some locational data to pull from. For example, all locations in the region of X may have price
				// Y.
				randomUnits, _ := randomPrice(priceLow, priceHigh)
				newRoute := esi.DerRoute{
					ExchangeKey: coordinationNodeInfo.PublicKey,
					FacilityKey: publicKey,
				}
				timeNow := timestamppb.Timestamp{
					Seconds: unixSeconds(),
					Nanos:   0,
				}
				newMoney := esi.Money{
					CurrencyCode: "USD",
					Units:        randomUnits,
					Nanos:        0,
				}
				newPriceComponents := esi.PriceComponents{
					ApparentEnergyPrice: &newMoney,
				}
				// Example test datum that could be sent to facilities.
				newDatum := esi.PriceDatum{
					Route:           &newRoute,
					Ts:              &timeNow,
					TimeUnit:        esi.TimeUnit_INSTANT,
					PriceComponents: &newPriceComponents,
				}
				_ = esi.ListPrices(coordinationNodeClient, &newDatum)
				log.WithFields(log.Fields{
					"dest":  newRoute.GetFacilityKey(),
					"price": newDatum.PriceComponents.ApparentEnergyPrice.Units,
				}).Info("Sent price datum")
			}
		}

		// Look at current offers.
		for uuid, offer := range priceMapOffers {
			// Actions specifically relating to the facility.
			if offer.Route.GetFacilityKey() == coordinationNodeInfo.GetPublicKey() {

				// If the offer has been accepted, then check to see if the time expected has passed.
				if priceMapOfferStatus[offer.OfferId.Uuid].Status == esi.PriceMapOfferStatus_ACCEPTED && offer.When.Seconds <= unixSeconds() {
					priceMapOfferStatus[offer.OfferId.Uuid].Status = esi.PriceMapOfferStatus_EXECUTING

					log.WithFields(log.Fields{
						"uuid": uuid,
					}).Info("Offer is executing")
				}
			}

			// If the offer is executing and has passed the time expected to execute, set it to complete.
			if priceMapOfferStatus[offer.OfferId.Uuid].Status == esi.PriceMapOfferStatus_EXECUTING {
				// The time when is in nanoseconds, and duration is in seconds.
				if (offer.When.Seconds + offer.PriceMap.Duration.Seconds) <= unixSeconds() {
					priceMapOfferStatus[offer.OfferId.Uuid].Status = esi.PriceMapOfferStatus_COMPLETED

					log.WithFields(log.Fields{
						"uuid": uuid,
					}).Info("Offer has completed")

					// Create a new feedback.
					newFeedback := esi.PriceMapOfferFeedback{
						Route:            offer.Route,
						OfferId:          offer.OfferId,
						ObligationStatus: 2,
					}

					// Get feedback from exchange.
					err := esi.GetPriceMapOfferFeedback(coordinationNodeClient, &newFeedback)
					if err != nil {
						log.Error(err.Error())
					}
				}
			}
		}

		// Do these actions at a regular interval.
		time.Sleep(time.Second * 20)
	}
}
