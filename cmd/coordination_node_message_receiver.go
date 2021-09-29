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
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// coordinationNodeMessageReceiver receives and returns any incoming coordination node messages.
func coordinationNodeMessageReceiver() {
	var formKey int // a simple number to increment form number.

	<-coordinationNodeClient.OnConnect.C
	log.WithFields(log.Fields{
		"publicKey": coordinationNodeInfo.GetPublicKey(),
		"name":      coordinationNodeInfo.GetName(),
	}).Info("Connection opened")

	message := &esi.CoordinationNodeMessage{}

	for {
		// Unmarshal the protocol buffer.
		msg := <-coordinationNodeClient.OnMessage.C
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			log.Error(err.Error())
		}

		// Case documentation located at api/esi/der_facility_service.go.
		//
		// Switch based upon the message type.
		switch x := message.Chunk.(type) {
		case *esi.CoordinationNodeMessage_SendKnownDerFacility:
			// If the node is not already stored, store it.
			_, present := knownCoordinationNodes[x.SendKnownDerFacility.GetPublicKey()]
			if !present {
				knownCoordinationNodes[x.SendKnownDerFacility.PublicKey] = x.SendKnownDerFacility
			}

			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info(fmt.Sprintf("Saved coordination node %s", x.SendKnownDerFacility.GetPublicKey()))

		case *esi.CoordinationNodeMessage_GetDerFacilityRegistrationForm:
			// Set the basic info.
			//
			// An example FormSetting - you can set whatever you want, and the facility will get a copy for you to then
			// evaluate as you wish.
			newFormSetting := esi.FormSetting{
				Key:         "0",
				Label:       "Do you wish to register?",
				Caption:     "",
				Placeholder: "Y",
			}
			newForm := esi.Form{
				LanguageCode: "en",
				Key:          strconv.Itoa(formKey),
				Settings:     []*esi.FormSetting{&newFormSetting},
			}
			newRoute := esi.DerRoute{
				FacilityKey: msg.Src,
				ExchangeKey: coordinationNodeInfo.GetPublicKey(),
			}
			newRegistrationForm := esi.DerFacilityRegistrationForm{
				Route: &newRoute,
				Form:  &newForm,
			}

			// Send the registration form.
			err = esi.SendDerFacilityRegistrationForm(coordinationNodeClient, &newRegistrationForm)
			if err != nil {
				log.Error(err.Error())
			}

			formKey += 1 // increment form key

			log.WithFields(log.Fields{
				"dest": msg.Src,
			}).Info("Sent registration form")

		case *esi.CoordinationNodeMessage_SendDerFacilityRegistrationForm:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received registration form")

			// If the form is not already stored, store it.
			_, present := receivedRegistrationForms[x.SendDerFacilityRegistrationForm.Route.GetExchangeKey()]
			if !present {
				receivedRegistrationForms[x.SendDerFacilityRegistrationForm.Route.GetExchangeKey()] = x.SendDerFacilityRegistrationForm
			}

		case *esi.CoordinationNodeMessage_SubmitDerFacilityRegistrationForm:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received registration form data")

			registration := esi.DerFacilityRegistration{
				Route: x.SubmitDerFacilityRegistrationForm.Route,
			}
			// If the user responded positively, then success.
			response := strings.ToLower(x.SubmitDerFacilityRegistrationForm.Data.Data["0"])
			if response == "y" || response == "yes" {
				registration.Success = true
			} else {
				registration.Success = false
			}

			// If successful, add it as a facility.
			if registration.Success {
				registeredFacilities[msg.Src] = true
			}

			err = esi.CompleteDerFacilityRegistration(coordinationNodeClient, &registration)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"dest":    msg.Src,
				"success": registration.GetSuccess(),
			}).Info("Sent completed registration form")

		case *esi.CoordinationNodeMessage_CompleteDerFacilityRegistration:
			if x.CompleteDerFacilityRegistration.GetSuccess() == true {
				registeredExchange = msg.Src
			}
			log.WithFields(log.Fields{
				"src":     msg.Src,
				"success": x.CompleteDerFacilityRegistration.GetSuccess(),
			}).Info("Received completed registration form")

			newRequest := esi.DerPowerParametersRequest{
				Route: x.CompleteDerFacilityRegistration.Route,
			}

			err = esi.GetPowerParameters(coordinationNodeClient, &newRequest)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"dest": msg.Src,
			}).Info("Getting power parameters")

		case *esi.CoordinationNodeMessage_GetPowerParameters:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Requested power parameters")

			// Send the power parameters associated with the service.
			//
			// At the moment, both nodes have the same power parameters set, so this doesn't really do anything. But
			// this shows that you can get the power parameters from another service, and having to set your own is
			// tedious for a demo.
			err = esi.SetPowerParameters(coordinationNodeClient, x.GetPowerParameters.Route.GetFacilityKey(), &powerParameters)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"dest": msg.Src,
			}).Info("Sent power parameters")

		case *esi.CoordinationNodeMessage_SetPowerParameters:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received power parameters")

			// Set your power parameters to the ones provided by the service.
			powerParameters = *x.SetPowerParameters

			log.WithFields(log.Fields{
				"src":   msg.Src,
				"param": x.SetPowerParameters,
			}).Info("Set power parameters")

		case *esi.CoordinationNodeMessage_ListPrices:
			log.WithFields(log.Fields{
				"src": msg.Src,
				"price": x.ListPrices.PriceComponents.ApparentEnergyPrice.Units,
			}).Info("Received price datum")

		case *esi.CoordinationNodeMessage_GetResourceCharacteristics:
			// Check to make sure that the source is the registered exchange.
			if registeredExchange == msg.Src {
				newRoute := esi.DerRoute{
					FacilityKey: coordinationNodeInfo.GetPublicKey(),
					ExchangeKey: msg.Src,
				}
				// TODO: fix
				newCharacteristics := resourceCharacteristics
				newCharacteristics.Route = &newRoute
				err := esi.SendResourceCharacteristics(coordinationNodeClient, &newCharacteristics)
				if err != nil {
					log.Error(err.Error())
				}

				log.WithFields(log.Fields{
					"dest": msg.Src,
				}).Info("Sent resource characteristics")
			}

		case *esi.CoordinationNodeMessage_SendResourceCharacteristics:
			// Check to make sure that the source is a registered facility.
			if _, ok := registeredFacilities[msg.Src]; ok {
				facilityCharacteristics[msg.Src] = x.SendResourceCharacteristics

				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Received resource characteristics")
			}

		case *esi.CoordinationNodeMessage_GetPriceMap:
			// Check to make sure that the source is the registered exchange.
			if registeredExchange == msg.Src {
				err = esi.SendPriceMap(coordinationNodeClient, x.GetPriceMap.Route.GetExchangeKey(), &priceMap)
				if err != nil {
					log.Error(err.Error())
				}

				log.WithFields(log.Fields{
					"dest": msg.Src,
				}).Info("Sent price map")
			}

		case *esi.CoordinationNodeMessage_SendPriceMap:
			// Check to make sure that the source is a registered facility.
			if _, ok := registeredFacilities[msg.Src]; ok {
				facilityPriceMaps[msg.Src] = x.SendPriceMap

				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Received price map")
			}

		case *esi.CoordinationNodeMessage_ProposePriceMapOffer:
			// Check to make sure that the source is the registered exchange or facility.
			_, ok := registeredFacilities[msg.Src]
			if registeredExchange == msg.Src || ok {
				log.Info("RECEIVED PROPOSE OFFER")
				if x.ProposePriceMapOffer.PriceMap.Price.ApparentEnergyPrice.Units < autoPrice.AlwaysBuyBelowPrice.Units {
					// If the offer is below our auto accept, just accept the offer.
					//
					// There is also a value for "AvoidBuyOverPrice", which could be used in a similar way in other
					// scenarios. In this demo, if the price is not lower than our auto accept, then it just goes to
					// evaluation.
					response := acceptOffer(x.ProposePriceMapOffer.Route, x.ProposePriceMapOffer.OfferId)
					err = esi.SendPriceMapOfferResponse(coordinationNodeClient, response)
					if err != nil {
						log.Error(err.Error())
					}

					log.WithFields(log.Fields{
						"src":  msg.Src,
						"auto": autoPrice.AlwaysBuyBelowPrice.Units,
					}).Info("Accepted price map due to auto buy")

					priceMapOffers[x.ProposePriceMapOffer.OfferId.Uuid] = x.ProposePriceMapOffer
					// Store the status of the offer.
					status := esi.PriceMapOfferStatus{
						Route:   x.ProposePriceMapOffer.Route,
						OfferId: x.ProposePriceMapOffer.OfferId,
						Status:  esi.PriceMapOfferStatus_ACCEPTED,
					}
					priceMapOfferStatus[x.ProposePriceMapOffer.OfferId.Uuid] = &status
				} else {
					priceMapOffers[x.ProposePriceMapOffer.OfferId.Uuid] = x.ProposePriceMapOffer

					log.WithFields(log.Fields{
						"src": msg.Src,
					}).Info("Received price map offer")

					// Store the status of the offer.
					status := esi.PriceMapOfferStatus{
						Route:   x.ProposePriceMapOffer.Route,
						OfferId: x.ProposePriceMapOffer.OfferId,
						Status:  esi.PriceMapOfferStatus_UNKNOWN,
					}
					priceMapOfferStatus[x.ProposePriceMapOffer.OfferId.Uuid] = &status
				}
			}

		case *esi.CoordinationNodeMessage_SendPriceMapOfferResponse:
			switch y := x.SendPriceMapOfferResponse.AcceptOneof.(type) {
			// Evaluate the contents of the response.
			case *esi.PriceMapOfferResponse_Accept:
				if y.Accept {
					// If the offer has been accepted, log the acceptance.
					log.WithFields(log.Fields{
						"src": msg.Src,
					}).Info("Price map accepted")

					// Store the status ACCEPTED.
					priceMapOfferStatus[x.SendPriceMapOfferResponse.OfferId.Uuid].Status = esi.PriceMapOfferStatus_ACCEPTED
				}
			case *esi.PriceMapOfferResponse_CounterOffer:
				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Counter offer received")

				// Store the previous offer as REJECTED.
				priceMapOfferStatus[x.SendPriceMapOfferResponse.PreviousOffer.Uuid].Status = esi.PriceMapOfferStatus_REJECTED

				// In the new offer, use the time specified by the previous offer.
				newOffer := esi.PriceMapOffer{
					Route:    x.SendPriceMapOfferResponse.Route,
					OfferId:  x.SendPriceMapOfferResponse.OfferId,
					When:     priceMapOffers[x.SendPriceMapOfferResponse.PreviousOffer.Uuid].When,
					PriceMap: x.SendPriceMapOfferResponse.GetCounterOffer(),
					//Node:     &esi.NodeType{Type: party},
					Node: x.SendPriceMapOfferResponse.Node,
				}
				// Store the new offer.
				priceMapOffers[x.SendPriceMapOfferResponse.OfferId.Uuid] = &newOffer

				// Store the status of the offer.
				status := esi.PriceMapOfferStatus{
					Route:   x.SendPriceMapOfferResponse.Route,
					OfferId: x.SendPriceMapOfferResponse.OfferId,
					Status:  esi.PriceMapOfferStatus_UNKNOWN,
				}
				priceMapOfferStatus[x.SendPriceMapOfferResponse.OfferId.Uuid] = &status

				if y.CounterOffer.Price.ApparentEnergyPrice.Units < autoPrice.AlwaysBuyBelowPrice.Units {
					// If it falls below the auto accept, then accept it.
					response := acceptOffer(x.SendPriceMapOfferResponse.Route, x.SendPriceMapOfferResponse.OfferId)
					err = esi.SendPriceMapOfferResponse(coordinationNodeClient, response)
					if err != nil {
						log.Error(err.Error())
					}

					priceMapOfferStatus[x.SendPriceMapOfferResponse.OfferId.Uuid].Status = esi.PriceMapOfferStatus_ACCEPTED

					log.WithFields(log.Fields{
						"src":  msg.Src,
						"auto": autoPrice.AlwaysBuyBelowPrice.Units,
					}).Info("Accepted price map due to auto buy")
				}
			}

		case *esi.CoordinationNodeMessage_GetPriceMapOfferFeedback:
			// As mentioned in coordination_node_input_receiver.go, this is merely a stub of what could be implemented.
			//
			// In a real situation, getting feedback on a response (either manually or automatically) is very powerful,
			// this is just to show the capability.
			if _, ok := registeredFacilities[msg.Src]; ok {
				log.WithFields(log.Fields{
					"src":   msg.Src,
					"claim": x.GetPriceMapOfferFeedback.ObligationStatus,
				}).Info("Received offer feedback")

				response := esi.PriceMapOfferFeedbackResponse{
					Route:    x.GetPriceMapOfferFeedback.Route,
					OfferId:  x.GetPriceMapOfferFeedback.OfferId,
					Accepted: true,
				}

				err := esi.ProvidePriceMapOfferFeedback(coordinationNodeClient, &response)
				if err != nil {
					log.Error(err.Error())
				}
			}

		case *esi.CoordinationNodeMessage_ProvidePriceMapOfferFeedback:
			log.WithFields(log.Fields{
				"src":   msg.Src,
				"claim": x.ProvidePriceMapOfferFeedback.Accepted,
			}).Info("Provide feedback response")
		}
	}
}

func acceptOffer(route *esi.DerRoute, offerId *esi.Uuid) *esi.PriceMapOfferResponse {
	accept := esi.PriceMapOfferResponse_Accept{
		Accept: true,
	}
	response := esi.PriceMapOfferResponse{
		Route:       route,
		OfferId:     offerId,
		AcceptOneof: &accept,
	}

	return &response
}
