package cmd

import (
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// facilityMessageReceiver receives and returns any incoming facility messages.
func facilityMessageReceiver() {
	var formKey float64 // a simple number to increment form number.

	<-facilityClient.OnConnect.C
	log.WithFields(log.Fields{
		"publicKey": facilityInfo.GetFacilityPublicKey(),
		"name":      facilityInfo.GetName(),
	}).Info("Connection opened")

	message := &esi.FacilityMessage{}

	// TODO: User created? Pass in as argument?
	// An example setting for form.
	formSetting := esi.FormSetting{
		Key:         "0",
		Label:       "Do you like apples?",
		Caption:     "",
		Placeholder: "Y",
	}
	// An example English language form.
	enForm := esi.Form{
		LanguageCode: "en",
		Key:          fmt.Sprintf("%f", formKey),
		Settings:     []*esi.FormSetting{&formSetting},
	}
	// An example English language registration form.
	registrationForm := esi.DerFacilityRegistrationForm{
		ProviderFacilityPublicKey: formatBinary(facilityClient.PubKey()),
		CustomerFacilityPublicKey: "", // fill in customer key when sending
		Form:                      &enForm,
	}

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
			_, present := knownFacilities[x.SendKnownDerFacility.GetFacilityPublicKey()]
			if !present {
				knownFacilities[x.SendKnownDerFacility.FacilityPublicKey] = x.SendKnownDerFacility
			}

			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info(fmt.Sprintf("Saved facility %s", x.SendKnownDerFacility.GetFacilityPublicKey()))

		case *esi.FacilityMessage_GetDerFacilityRegistrationForm:
			// Set the basic info.
			registrationForm.CustomerFacilityPublicKey = msg.Src
			registrationForm.ProviderFacilityPublicKey = facilityInfo.GetFacilityPublicKey()

			// Send the registration form.
			err = esi.SendDerFacilityRegistrationForm(facilityClient, registrationForm)
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
			_, present := receivedRegistrationForms[x.SendDerFacilityRegistrationForm.GetProviderFacilityPublicKey()]
			if !present {
				receivedRegistrationForms[x.SendDerFacilityRegistrationForm.GetProviderFacilityPublicKey()] = x.SendDerFacilityRegistrationForm
			}

		case *esi.FacilityMessage_SubmitDerFacilityRegistrationForm:
			// TODO: Fill in registration form.
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received registration form data")

			// INFO: In this case, automatically accept registration forms.
			registration := esi.DerFacilityRegistration{
				Route: x.SubmitDerFacilityRegistrationForm.Route,
				Success: true,
			}

			err = esi.CompleteDerFacilityRegistration(facilityClient, registration)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"end": msg.Src,
				"success": registration.GetSuccess(),
			}).Info("Sent completed registration form")

		case *esi.FacilityMessage_CompleteDerFacilityRegistration:
			log.WithFields(log.Fields{
				"src": msg.Src,
				"success": x.CompleteDerFacilityRegistration.GetSuccess(),
			}).Info("Received completed registration form")
		}
	}
}
