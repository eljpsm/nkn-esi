package cmd

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// receivedRegistrationForms is a list of registration forms the user can then fill out.
var receivedRegistrationForms = []esi.DerFacilityRegistrationForm{}

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
		Label:       "Do you like apples?: ",
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
			isInKnownFacilities := false
			for _, v := range knownFacilities {
				if v.GetFacilityPublicKey() == x.SendKnownDerFacility.GetFacilityPublicKey() {
					isInKnownFacilities = true
				}
			}
			if isInKnownFacilities == false {
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

			isInForms := false
			for _, v := range receivedRegistrationForms {
				if v.GetProviderFacilityPublicKey() == x.SendDerFacilityRegistrationForm.GetProviderFacilityPublicKey() {
					isInForms = true
				}
			}
			if isInForms == false {
				receivedRegistrationForms = append(receivedRegistrationForms, *x.SendDerFacilityRegistrationForm)
			}

		case *esi.FacilityMessage_SubmitDerFacilityRegistrationForm:
			// TODO: Fill in registration form.
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received registration form data")

			route := esi.DerRoute{
				BuyKey:  facilityInfo.GetFacilityPublicKey(),
				SellKey: msg.Src,
			}
			registration := esi.DerFacilityRegistration{
				Route: &route,
			}

			err = esi.CompleteDerFacilityRegistration(facilityClient, registration)
			if err != nil {
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"end": msg.Src,
			}).Info("Sent completed registration form")
			log.WithFields(log.Fields{
				"end": msg.Src,
			}).Info("Permission granted")

		case *esi.FacilityMessage_CompleteDerFacilityRegistration:
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Received completed registration form")
			log.WithFields(log.Fields{
				"end": msg.Src,
			}).Info("Granted permission")
		}
	}
}

// facilityInputReceiver receives and returns any facility inputs.
func facilityInputReceiver() {
	shell := ishell.New()
	<-facilityClient.OnConnect.C
	shell.Printf("Connection opened on facility '%s'\n", infoMsgColorFunc(facilityInfo.GetName()))

	shell.AddCmd(&ishell.Cmd{
		Name: "public",
		Help: "print public key",
		Func: func(c *ishell.Context) {
			fmt.Println(facilityInfo.GetFacilityPublicKey())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "facilities",
		Help: "print known facilities",
		Func: func(c *ishell.Context) {
			if len(knownFacilities) > 0 {
				for _, v := range knownFacilities {
					shell.Printf("\nName: %s\nCountry: %s\nPublic Key: %s\n", v.GetName(), v.Location.GetCountry(), noteMsgColorFunc(v.GetFacilityPublicKey()))
				}
				shell.Println()
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "signup",
		Help: "sign up to a registry",
		Func: func(c *ishell.Context) {
			c.Print("Registry Public Key: ")

			err := esi.SignupRegistry(facilityClient, c.ReadLine(), facilityInfo)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "query",
		Help: "query registry for facilities by location",
		Func: func(c *ishell.Context) {
			// TODO: Ask for more than just country.
			// TODO: Evaluate settings.
			c.Print("Registry Public Key: ")
			registryPublicKey := c.ReadLine()
			c.Print("Country: ")
			country := c.ReadLine()
			newLocation := esi.Location{
				Country: country,
			}

			request := esi.DerFacilityExchangeRequest{Location: &newLocation}
			err := esi.QueryDerFacilities(facilityClient, registryPublicKey, request)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "request",
		Help: "request registration form from facility",
		Func: func(c *ishell.Context) {
			c.Print("Facility Public Key: ")
			facilityPublicKey := c.ReadLine()
			c.Print("Language Code: ")
			languageCode := c.ReadLine()

			request := esi.DerFacilityRegistrationFormRequest{
				FacilityPublicKey: facilityPublicKey,
				LanguageCode:      languageCode,
			}

			err := esi.GetDerFacilityRegistrationForm(facilityClient, request)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "forms",
		Help: "print forms to be signed",
		Func: func(c *ishell.Context) {
			if len(receivedRegistrationForms) > 0 {
				for _, v := range receivedRegistrationForms {
					shell.Printf("\nProvider Public Key: %s\n", noteMsgColorFunc(v.GetProviderFacilityPublicKey()))
				}
				shell.Println()
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "sign",
		Help: "fill in a received registration form",
		Func: func(c *ishell.Context) {
			c.Print("Facility Public Key: ")
			facilityPublicKey := c.ReadLine()

			signed := false
			for i, v := range receivedRegistrationForms {
				if v.GetProviderFacilityPublicKey() == facilityPublicKey {
					// TODO: nonce
					data := esi.DerFacilityRegistrationFormData{
						CustomerFacilityPublicKey: facilityInfo.GetFacilityPublicKey(),
					}

					err := esi.SubmitDerFacilityRegistrationForm(facilityClient, data)
					if err != nil {
						log.Error(err.Error())
					}

					signed = true

					// Remove form from list.
					receivedRegistrationForms[i] = receivedRegistrationForms[len(receivedRegistrationForms)-1]
					receivedRegistrationForms[len(receivedRegistrationForms)-1] = esi.DerFacilityRegistrationForm{}
					receivedRegistrationForms = receivedRegistrationForms[:len(receivedRegistrationForms)-1]

					log.WithFields(log.Fields{
						"end": v.GetProviderFacilityPublicKey(),
					}).Info("Sent registration form")
				}
			}
			if !signed {
				shell.Printf("No form found with public key '%s`\n", facilityPublicKey)
			}
		},
	})

	shell.Run()
}
