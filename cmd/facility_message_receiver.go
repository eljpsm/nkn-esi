package cmd

import (
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// facilityMessageReceiver receives and returns any incoming facility messages.
func facilityMessageReceiver() {
	var formKey int // a simple number to increment form number.

	<-facilityClient.OnConnect.C
	log.WithFields(log.Fields{
		"publicKey": facilityInfo.GetPublicKey(),
		"name":      facilityInfo.GetName(),
	}).Info("Connection opened")

	message := &esi.FacilityMessage{}

	for {
		// Unmarshal the protocol buffer.
		msg := <-facilityClient.OnMessage.C
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			log.Error(err.Error())
		}

		// Case documentation located at api/esi/deer_facility_service.go.
		//
		// Switch based upon the message type.
		switch x := message.Chunk.(type) {
		case *esi.FacilityMessage_SendKnownDerFacility:
			// If the facility is not already stored, store it.
			_, present := knownFacilities[x.SendKnownDerFacility.GetPublicKey()]
			if !present {
				knownFacilities[x.SendKnownDerFacility.PublicKey] = x.SendKnownDerFacility
			}

			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info(fmt.Sprintf("Saved facility %s", x.SendKnownDerFacility.GetPublicKey()))

		case *esi.FacilityMessage_GetDerFacilityRegistrationForm:
			// Set the basic info.
			//
			// An example FormSetting - you can set whatever you want, and the producer will get a copy for you to then
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
			newRegistrationForm := esi.DerFacilityRegistrationForm{
				ProducerKey: facilityInfo.GetPublicKey(),
				CustomerKey: msg.Src,
				Form:        &newForm,
			}

			// Send the registration form.
			err = esi.SendDerFacilityRegistrationForm(facilityClient, newRegistrationForm)
			if err != nil {
				log.Error(err.Error())
			}

			formKey += 1 // increment form key

			log.WithFields(log.Fields{
				"end": msg.Src,
			}).Info("Sent registration form")

		case *esi.FacilityMessage_SendDerFacilityRegistrationForm:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received registration form")

			// If the form is not already stored, store it.
			_, present := receivedRegistrationForms[x.SendDerFacilityRegistrationForm.GetProducerKey()]
			if !present {
				receivedRegistrationForms[x.SendDerFacilityRegistrationForm.GetProducerKey()] = x.SendDerFacilityRegistrationForm
			}

		case *esi.FacilityMessage_SubmitDerFacilityRegistrationForm:
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

			// If successful, add it as a consumer facility with an empty price map.
			if registration.Success {
				producerFacilities[msg.Src] = true
			}

			err = esi.CompleteDerFacilityRegistration(facilityClient, registration)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"end":     msg.Src,
				"success": registration.GetSuccess(),
			}).Info("Sent completed registration form")

		case *esi.FacilityMessage_CompleteDerFacilityRegistration:
			if x.CompleteDerFacilityRegistration.GetSuccess() == true {
				customerFacility = msg.Src
			}
			log.WithFields(log.Fields{
				"src":     msg.Src,
				"success": x.CompleteDerFacilityRegistration.GetSuccess(),
			}).Info("Received completed registration form")

		case *esi.FacilityMessage_GetResourceCharacteristics:
			// Check to make sure that the source is a registered customer.
			if customerFacility != "" {
				newRoute := esi.DerRoute{
					CustomerKey: msg.Src,
					ProducerKey: facilityInfo.GetPublicKey(),
				}
				newCharacteristics := resourceCharacteristics
				newCharacteristics.Route = &newRoute
				err := esi.SendResourceCharacteristics(facilityClient, newCharacteristics)
				if err != nil {
					log.Error(err.Error())
				}

				log.WithFields(log.Fields{
					"end": msg.Src,
				}).Info("Sent resource characteristics")
			}

		case *esi.FacilityMessage_SendResourceCharacteristics:
			// Check to make sure that the source is a registered producer.
			if producerFacilities[msg.Src] == true {
				producerCharacteristics[msg.Src] = x.SendResourceCharacteristics

				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Received resource characteristics")
			}

		case *esi.FacilityMessage_GetPriceMap:
			// Check to make sure that the source is a registered customer.
			if customerFacility == msg.Src {
				err = esi.SendPriceMap(facilityClient, x.GetPriceMap.Route.GetCustomerKey(), priceMap)
				if err != nil {
					log.Error(err.Error())
				}

				log.WithFields(log.Fields{
					"end": msg.Src,
				}).Info("Sent price map")
			}

		case *esi.FacilityMessage_SendPriceMap:
			// Check to make sure that the source is a registered producer.
			if producerFacilities[msg.Src] == true {
				producerPriceMaps[msg.Src] = x.SendPriceMap

				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Received price map")
			}

		case *esi.FacilityMessage_ProposePriceMapOffer:
			// Check to make sure that the source is a registered customer.
			if customerFacility == msg.Src {
				if x.ProposePriceMapOffer.PriceMap.Price.ApparentEnergyPrice.Units < autoPrice.AlwaysBuyBelowPrice.Units {
					// If the offer is below our auto accept, just accept the offer.
					//
					// There is also a value for "AvoidBuyOverPrice", which could be used in a similar way in other
					// scenarios. In this demo, if the price is not lower than our auto accept, then it just goes to
					// evaluation.
					accept := esi.PriceMapOfferResponse_Accept{
						Accept: true,
					}
					response := esi.PriceMapOfferResponse{
						Route:       x.ProposePriceMapOffer.Route,
						OfferId:     x.ProposePriceMapOffer.OfferId,
						AcceptOneof: &accept,
					}
					err = esi.SendPriceMapOfferResponse(facilityClient, response)
					if err != nil {
						log.Error(err.Error())
					}

					log.WithFields(log.Fields{
						"src": msg.Src,
						"auto": autoPrice.AlwaysBuyBelowPrice.Units,
					}).Info("Accepted price map due to auto buy")
				} else {
					priceMapOffers[x.ProposePriceMapOffer.OfferId.Uuid] = x.ProposePriceMapOffer
				}

				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Received price map offer")
			}

		case *esi.FacilityMessage_ProvidePriceMapOfferFeedback:
			if producerFacilities[msg.Src] == true {
			}
		}
	}
}
