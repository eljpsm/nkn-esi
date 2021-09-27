package cmd

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/ptypes/duration"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var (
	// priceMapOffers is the currently stored price map offers.
	priceMapOffers = make(map[string]*esi.PriceMapOffer)

	// priceMapOfferFeedbacks is the currently stored price map offer feedbacks.
	priceMapOfferFeedbacks = make(map[string]*esi.PriceMapOfferFeedback)

	// priceMaps is the currently stored price map per negotiation.
	priceMaps = make(map[string]*esi.PriceMap)

	// receivedRegistrationForms is a list of registration forms the user can then fill out.
	receivedRegistrationForms = make(map[string]*esi.DerFacilityRegistrationForm)
)

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

			v, present := receivedRegistrationForms[facilityPublicKey]
			if present {
				shell.Println() // gap from input
				fmt.Println(v.GetCustomerFacilityPublicKey())
				fmt.Println(v.GetProviderFacilityPublicKey())

				// TODO: nonce
				// Contains the results of key -> response.
				results := make(map[string]string)
				route := esi.DerRoute{
					BuyKey:  v.GetCustomerFacilityPublicKey(),
					SellKey: v.GetProviderFacilityPublicKey(),
				}

				formData := esi.FormData{
					Data: results,
				}
				// Contains the full form data.
				registrationFormData := esi.DerFacilityRegistrationFormData{
					Route: &route,
					Data:  &formData,
				}

				for _, v := range v.Form.Settings {
					// For all the settings, print the desired setting, get an input and then store it in the
					// results.
					shell.Printf("%s. %s [%s]: ", v.GetKey(), v.GetLabel(), v.GetPlaceholder())
					result := c.ReadLine()

					// If input is not given, then use the placeholder value.
					if v.GetPlaceholder() != "" {
						if result == "" {
							result = v.GetPlaceholder()
						}
					}

					results[v.Key] = result
				}

				// Submit the registration form.
				err := esi.SubmitDerFacilityRegistrationForm(facilityClient, registrationFormData)
				if err != nil {
					log.Error(err.Error())
				}

				// Remove form from the map.
				delete(receivedRegistrationForms, v.GetProviderFacilityPublicKey())

				log.WithFields(log.Fields{
					"end": v.GetProviderFacilityPublicKey(),
				}).Info("Sent registration form")

				shell.Printf("\nForm has been submitted to %s\n", registrationFormData.Route.GetSellKey())

			} else {
				shell.Printf("no form found with public key '%s`\n", facilityPublicKey)
				return
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "price-map",
		Help: "create a price map",
		Func: func(c *ishell.Context) {

			// Create newPowerComponents.
			shell.Print("Real Power: ")
			realPowerString := c.ReadLine()
			realPower, err := strconv.Atoi(realPowerString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Reactive Power: ")
			reactivePowerString := c.ReadLine()
			reactivePower, err := strconv.Atoi(reactivePowerString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newPowerComponents := esi.PowerComponents{
				RealPower:     int64(realPower),
				ReactivePower: int64(reactivePower),
			}

			// Create newDuration.
			shell.Print("Expected Duration Seconds: ")
			durationSecondsString := c.ReadLine()
			durationSeconds, err := strconv.Atoi(durationSecondsString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Expected Duration Nanos: ")
			durationNanosString := c.ReadLine()
			durationNanos, err := strconv.Atoi(durationNanosString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newDuration := duration.Duration{
				Seconds: int64(durationSeconds),
				Nanos:   int32(durationNanos),
			}

			// Create newDurationRange.
			shell.Print("Expected Minimum Duration Seconds: ")
			expectedMinSecondsString := c.ReadLine()
			expectedMinSeconds, err := strconv.Atoi(expectedMinSecondsString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Expected Minimum Duration Nanos: ")
			expectedMinNanosString := c.ReadLine()
			expectedMinNanos, err := strconv.Atoi(expectedMinNanosString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Expected Maximum Duration Seconds: ")
			expectedMaxSecondsString := c.ReadLine()
			expectedMaxSeconds, err := strconv.Atoi(expectedMaxSecondsString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Expected Maximum Duration Nanos: ")
			expectedMaxNanosString := c.ReadLine()
			expectedMaxNanos, err := strconv.Atoi(expectedMaxNanosString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newMinDuration := duration.Duration{
				Seconds: int64(expectedMinSeconds),
				Nanos:   int32(expectedMinNanos),
			}
			newMaxDuration := duration.Duration{
				Seconds: int64(expectedMaxSeconds),
				Nanos:   int32(expectedMaxNanos),
			}
			newDurationRange := esi.DurationRange{
				Min: &newMinDuration,
				Max: &newMaxDuration,
			}

			// Create newPriceComponents.
			shell.Print("Currency Code: ")
			currencyCode := c.ReadLine()
			shell.Print("Units: ")
			unitsString := c.ReadLine()
			units, err := strconv.Atoi(unitsString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Nanos: ")
			nanosString := c.ReadLine()
			nanos, err := strconv.Atoi(nanosString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newMoney := esi.Money{
				CurrencyCode: currencyCode,
				Units:        int64(units),
				Nanos:        int32(nanos),
			}
			newPriceComponents := esi.PriceComponents{
				ApparentEnergyPrice: &newMoney,
			}

			newPriceMap := esi.PriceMap{
				PowerComponents: &newPowerComponents,
				Duration:        &newDuration,
				ResponseTime:    &newDurationRange,
				Price:           &newPriceComponents,
			}

			fmt.Println(newPriceMap)
		},
	})

	shell.Run()
}
