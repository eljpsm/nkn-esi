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
		msg := <-facilityClient.OnMessage.C
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			log.Error(err.Error())
		}

		// Case documentation located at api/esi/deer_facility_service.go.
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
				Form:                      &newForm,
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
				customerFacilities[msg.Src] = true
			}
			log.WithFields(log.Fields{
				"src":     msg.Src,
				"success": x.CompleteDerFacilityRegistration.GetSuccess(),
			}).Info("Received completed registration form")

		case *esi.FacilityMessage_GetResourceCharacteristics:
			// Check to make sure that the source is a registered customer.
			if customerFacilities[msg.Src] == true {
				newRoute := esi.DerRoute{
					CustomerKey: msg.Src,
					ProducerKey: facilityInfo.GetPublicKey(),
				}
				newCharacteristics := resourceCharacteristics
				newCharacteristics.Route = &newRoute
				esi.SendResourceCharacteristics(facilityClient, newCharacteristics)

				log.WithFields(log.Fields{
					"end":     msg.Src,
				}).Info("Sent resource characteristics")
			}

		case *esi.FacilityMessage_SendResourceCharacteristics:
			// Check to make sure that the source is a registered producer.
			if producerFacilities[msg.Src] == true {
				producerCharacteristics[msg.Src] = x.SendResourceCharacteristics

				log.WithFields(log.Fields{
					"src":     msg.Src,
				}).Info("Received resource characteristics")
			}

		case *esi.FacilityMessage_GetPriceMap:
			// Check to make sure that the source is a registered customer.
			if customerFacilities[msg.Src] == true {
				esi.SendPriceMap(facilityClient, x.GetPriceMap.Route.GetCustomerKey(), priceMap)

				log.WithFields(log.Fields{
					"end":     msg.Src,
				}).Info("Sent price map")
			}

		case *esi.FacilityMessage_SendPriceMap:
			// Check to make sure that the source is a registered producer.
			if producerFacilities[msg.Src] == true {
				producerPriceMaps[msg.Src] = x.SendPriceMap

				log.WithFields(log.Fields{
					"src":     msg.Src,
				}).Info("Received price map")
			}

		case *esi.FacilityMessage_ProposePriceMapOffer:
			if customerFacilities[msg.Src] == true {
				// TODO: autobuy
				// TODO: what does ignore look like?
				customerPriceMapOffers[msg.Src] = x.ProposePriceMapOffer

				log.WithFields(log.Fields{
					"src":     msg.Src,
				}).Info("Received price map offer")
			}
		}
	}
}
