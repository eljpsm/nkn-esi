package cmd

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/ptypes/duration"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
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
			fmt.Println(facilityInfo.GetPublicKey())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "facilities",
		Help: "print known facilities",
		Func: func(c *ishell.Context) {
			if len(knownFacilities) > 0 {
				for _, facility := range knownFacilities {
					shell.Printf("\nName: %s\nCountry: %s\nPublic Key: %s\n", facility.GetName(), facility.Location.GetCountry(), noteMsgColorFunc(facility.GetPublicKey()))
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
				PublicKey:    facilityPublicKey,
				LanguageCode: languageCode,
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
					shell.Printf("\nProducer Public Key: %s\n", noteMsgColorFunc(v.GetProducerKey()))
				}
				shell.Println()
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "register",
		Help: "fill in a received registration form",
		Func: func(c *ishell.Context) {
			c.Print("Facility Public Key: ")
			facilityPublicKey := c.ReadLine()

			form, present := receivedRegistrationForms[facilityPublicKey]
			if present {
				shell.Println() // gap from input

				// TODO: nonce
				// Contains the results of key -> response.
				results := make(map[string]string)
				route := esi.DerRoute{
					CustomerKey: form.GetCustomerKey(),
					ProducerKey: form.GetProducerKey(),
				}

				formData := esi.FormData{
					Data: results,
				}
				// Contains the full form data.
				registrationFormData := esi.DerFacilityRegistrationFormData{
					Route: &route,
					Data:  &formData,
				}

				for _, setting := range form.Form.Settings {
					// For all the settings, print the desired setting, get an input and then store it in the
					// results.
					shell.Printf("%s. %s [%s]: ", setting.GetKey(), setting.GetLabel(), setting.GetPlaceholder())
					result := c.ReadLine()
					// If input is not given, then use the placeholder value.
					if setting.GetPlaceholder() != "" {
						if result == "" {
							result = setting.GetPlaceholder()
						}
					}

					results[setting.Key] = result
				}

				// Submit the registration form.
				err := esi.SubmitDerFacilityRegistrationForm(facilityClient, registrationFormData)
				if err != nil {
					log.Error(err.Error())
				}

				// Remove form from the map.
				delete(receivedRegistrationForms, form.GetProducerKey())

				log.WithFields(log.Fields{
					"end": form.GetProducerKey(),
				}).Info("Sent registration form")

				shell.Printf("\nForm has been submitted to %s\n", registrationFormData.Route.GetProducerKey())

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

			priceMap = newPriceMap
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "peek",
		Help: "view price map",
		Func: func(c *ishell.Context) {
			fmt.Println(priceMap)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "peers",
		Help: "show any registered customers or producers",
		Func: func(c *ishell.Context) {
			if len(customerFacilities) > 0 {
				shell.Println("\nCUSTOMERS")
				for k := range customerFacilities {
					shell.Printf("Facility Public Key: %s\n", k)
				}
			}
			if len(producerFacilities) > 0 {
				shell.Println("\nPRODUCERS")
				// TODO: add placeholders?
				// TODO: better formatting
				for k := range producerFacilities {
					shell.Printf("Facility Public Key: %s\n", k)
					if producerCharacteristics[k] != nil {
						shell.Println(producerCharacteristics[k])

					}
					if producerPriceMaps[k] != nil {
						shell.Println(producerPriceMaps[k])
					}
				}
				shell.Println()
			} else {
				shell.Println()
			}
		},
	})

	// TODO: add show characteristics
	shell.AddCmd(&ishell.Cmd{
		Name: "characteristics",
		Help: "create facility characteristics",
		Func: func(c *ishell.Context) {
			shell.Print("Max Load Power: ")
			loadPowerMaxString := c.ReadLine()
			loadPowerMax, err := strconv.Atoi(loadPowerMaxString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Load Power Factor: ")
			loadPowerFactorString := c.ReadLine()
			loadPowerFactor, err := strconv.Atoi(loadPowerFactorString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Max Supply Power: ")
			supplyPowerMaxString := c.ReadLine()
			supplyPowerMax, err := strconv.Atoi(supplyPowerMaxString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Supply Power Factor: ")
			supplyPowerFactorString := c.ReadLine()
			supplyPowerFactor, err := strconv.Atoi(supplyPowerFactorString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			shell.Print("Storage Energy Capacity: ")
			storageEnergyCapacityString := c.ReadLine()
			storageEnergyCapacity, err := strconv.Atoi(storageEnergyCapacityString)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			resourceCharacteristics.LoadPowerMax = uint64(loadPowerMax)
			resourceCharacteristics.LoadPowerFactor = float32(loadPowerFactor)
			resourceCharacteristics.SupplyPowerMax = uint64(supplyPowerMax)
			resourceCharacteristics.SupplyPowerFactor = float32(supplyPowerFactor)
			resourceCharacteristics.StorageEnergyCapacity = uint64(storageEnergyCapacity)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "get",
		Help: "get characteristics and price map of facility",
		Func: func(c *ishell.Context) {
			shell.Print("Producer Public Key: ")
			publicKey := c.ReadLine()

			if !producerFacilities[publicKey] {
				// TODO: better message, error?
				shell.Println("no key")
				return
			}

			newRoute := esi.DerRoute{
				CustomerKey: facilityInfo.GetPublicKey(),
				ProducerKey: publicKey,
			}
			newCharacteristicsRequest := esi.DerResourceCharacteristicsRequest{
				Route: &newRoute,
			}
			newPriceMapRequest := esi.DerPriceMapRequest{
				Route: &newRoute,
			}

			esi.GetResourceCharacteristics(facilityClient, newCharacteristicsRequest)
			esi.GetPriceMap(facilityClient, newPriceMapRequest)
		},
	})
	
	shell.AddCmd(&ishell.Cmd{
		Name: "propose",
		Help: "propose a price map offer to a producer",
		Func: func(c *ishell.Context) {
			shell.Print("Producer Public Key: ")
			publicKey := c.ReadLine()
			if !producerFacilities[publicKey] {
				// TODO: better message, error?
				// TODO: find way to combine with creating a new price map (input)
				shell.Println("no key")
				return
			}

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

			newRoute := esi.DerRoute{
				ProducerKey: publicKey,
				CustomerKey: facilityInfo.GetPublicKey(),
			}
			newPriceMap := esi.PriceMap{
				PowerComponents: &newPowerComponents,
				Duration:        &newDuration,
				ResponseTime:    &newDurationRange,
				Price:           &newPriceComponents,
			}

			newUuid := esi.Uuid{
				Hi: uuidHigh,
				Lo: uuidLow,
			}

			newTimeStamp := timestamppb.Timestamp{
				Seconds: unixSeconds(),
				Nanos: 0,
			}

			newPriceMapOffer := esi.PriceMapOffer{
				Route: &newRoute,
				OfferId: &newUuid,
				When: &newTimeStamp,
				PriceMap: &newPriceMap,
			}

			esi.ProposePriceMapOffer(facilityClient, newPriceMapOffer)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "offers",
		Help: "view pending offers",
		Func: func(c *ishell.Context) {
			for _, v := range customerPriceMapOffers {
				shell.Printf("\nCustomer Public Key: %s\nOffer: %s\n", noteMsgColorFunc(v.Route.GetProducerKey(), v.PriceMap))
			}
			shell.Println()
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "feedback",
		Help: "get feedback on a price map offer",
		Func: func(c *ishell.Context) {
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "evaluate",
		Help: "evaluate an offer and give feedback",
		Func: func(c *ishell.Context) {
			choices := []string{}
			for k, _ := range customerPriceMapOffers {
				choices = append(choices, k)
			}
			publicKey := c.MultiChoice(choices, "Select customer key to evaluate")
			shell.Println(publicKey)
			fmt.Printf("\n%s\n\n", customerPriceMapOffers[string(publicKey)])
		},
	})

	shell.Run()
}
