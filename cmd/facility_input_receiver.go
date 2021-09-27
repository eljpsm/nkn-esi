package cmd

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	uuid2 "github.com/gofrs/uuid"
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

	facilityInfoShellCmd := &ishell.Cmd{
		Name: "info",
		Help: "print facility information",
	}
	shell.AddCmd(facilityInfoShellCmd)
	facilityInfoShellCmd.AddCmd(&ishell.Cmd{
		Name: "public",
		Help: "print public key",
		Func: func(c *ishell.Context) {
			shell.Printf("%s\n", infoMsgColorFunc(facilityInfo.GetPublicKey()))
		},
	})

	facilityRegistryShellCmd := &ishell.Cmd{
		Name: "registry",
		Help: "manage registry functionality",
	}
	shell.AddCmd(facilityRegistryShellCmd)
	facilityRegistryShellCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "print known facilities received from registry",
		Func: func(c *ishell.Context) {
			for _, facility := range knownFacilities {
				// Print any information - currently only name, country, and public key.
				shell.Printf("\n%s %s\n%s %s\n%s %s\n",
					boldMsgColorFunc("Name:"),
					facility.GetName(),
					boldMsgColorFunc("Country:"),
					facility.Location.GetCountry(),
					boldMsgColorFunc("Public Key:"),
					noteMsgColorFunc(facility.GetPublicKey()))
			}
			shell.Println()
		},
	})
	facilityRegistryShellCmd.AddCmd(&ishell.Cmd{
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
	facilityRegistryShellCmd.AddCmd(&ishell.Cmd{
		Name: "query",
		Help: "query registry for facilities by location",
		Func: func(c *ishell.Context) {
			c.Print("Registry Public Key: ")
			registryPublicKey := c.ReadLine()
			c.Print("Country: ")
			country := c.ReadLine()

			// You can query based upon any setting that DerFacilityExchangeRequest takes.
			//
			// In this demo, you can select a COUNTRY to query based upon.
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
		Name: "peers",
		Help: "show any registered customers or producers",
		Func: func(c *ishell.Context) {
			if len(customerFacilities) > 0 {
				shell.Printf("\n%s\n", boldMsgColorFunc("CUSTOMERS"))
				for k := range customerFacilities {
					shell.Printf("%s %s\n",
						boldMsgColorFunc("Public Key:"),
						noteMsgColorFunc(k))
				}
			}
			if len(producerFacilities) > 0 {
				shell.Printf("\n%s\n", boldMsgColorFunc("PRODUCERS"))
				for k := range producerFacilities {
					shell.Printf("%s %s\n",
						boldMsgColorFunc("Public Key:"),
						noteMsgColorFunc(k))

					// TODO: pretty printing
					if producerCharacteristics[k] != nil {
						shell.Println(producerCharacteristics[k])
					}
					// TODO: pretty printing
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

	facilityProducerShellCmd := &ishell.Cmd{
		Name: "producer",
		Help: "manage producer functionality",
	}
	shell.AddCmd(facilityProducerShellCmd)
	facilityProducerShellCmd.AddCmd(&ishell.Cmd{
		Name: "request",
		Help: "request registration form from facility",
		Func: func(c *ishell.Context) {
			c.Print("Facility Public Key: ")
			facilityPublicKey := c.ReadLine()
			c.Print("Language Code: ")
			languageCode := c.ReadLine()

			// When creating a request, you can specify a language code.
			//
			// In this demo, the only language code that is used (and sent) is "en" for English.
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
	facilityProducerShellCmd.AddCmd(&ishell.Cmd{
		Name: "forms",
		Help: "print forms to be signed",
		Func: func(c *ishell.Context) {
			for _, v := range receivedRegistrationForms {
				shell.Printf("\n%s %s\n",
					noteMsgColorFunc("Public Key:"),
					noteMsgColorFunc(v.GetProducerKey()))
			}
			shell.Println()
		},
	})
	facilityProducerShellCmd.AddCmd(&ishell.Cmd{
		Name: "register",
		Help: "fill in a received registration form",
		Func: func(c *ishell.Context) {
			shell.Print("Facility Public Key: ")
			publicKey := c.ReadLine()

			form, present := receivedRegistrationForms[publicKey]
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
					//
					// This placeholder value given by DerFacilityRegistrationFormData is useful for any number of
					// situations in which user input could be either option or unnecessary.
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
				shell.Printf("no form found with public key '%s`\n", publicKey)
				return
			}
		},
	})

	facilityPriceMapShellCmd := &ishell.Cmd{
		Name: "price-map",
		Help: "manage local price maps",
	}
	shell.AddCmd(facilityPriceMapShellCmd)
	facilityPriceMapShellCmd.AddCmd(&ishell.Cmd{
		Name: "peek",
		Help: "view local price map",
		Func: func(c *ishell.Context) {
			fmt.Println(isPriceMapAccepted)
			fmt.Println(priceMap)
		},
	})
	facilityPriceMapShellCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "create a local price map",
		Func: func(c *ishell.Context) {
			createdPriceMap, err := newPriceMap(shell, c)
			if err != nil {
				shell.Println(err.Error())
				return
			}

			priceMap = createdPriceMap
		},
	})

	facilityCharacteristicsShellCmd := &ishell.Cmd{
		Name: "characteristics",
		Help: "manage local characteristics",
	}
	shell.AddCmd(facilityCharacteristicsShellCmd)
	facilityCharacteristicsShellCmd.AddCmd(&ishell.Cmd{
		Name: "peek",
		Help: "view local characteristics",
		Func: func(c *ishell.Context) {
			fmt.Println(resourceCharacteristics)
		},
	})
	facilityCharacteristicsShellCmd.AddCmd(&ishell.Cmd{
		Name: "create",
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

			// Set the local characteristics to user input.
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
				shell.Printf("no producer with public key: '%s\n'", publicKey)
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

			err := esi.GetResourceCharacteristics(facilityClient, newCharacteristicsRequest)
			if err != nil {
				shell.Println(err.Error())
			}
			err = esi.GetPriceMap(facilityClient, newPriceMapRequest)
			if err != nil {
				shell.Println(err.Error())
			}
		},
	})

	facilityCustomerShellCmd := &ishell.Cmd{
		Name: "customer",
		Help: "manage customer functionality",
	}
	shell.AddCmd(facilityCustomerShellCmd)
	facilityCustomerShellCmd.AddCmd(&ishell.Cmd{
		Name: "propose",
		Help: "propose a price map offer to a producer",
		Func: func(c *ishell.Context) {
			shell.Print("Producer Public Key: ")
			publicKey := c.ReadLine()
			if !producerFacilities[publicKey] {
				shell.Printf("no producer with public key: '%s'\n", publicKey)
				return
			}

			createdPriceMap, err := newPriceMap(shell, c)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newRoute := esi.DerRoute{
				ProducerKey: publicKey,
				CustomerKey: facilityInfo.GetPublicKey(),
			}
			// Each new offer can take a unique UUID for storage.
			//
			// In this demo, public keys are used instead - but depending on use case, this could be more useful.
			newUuid := esi.Uuid{
				Hi: uuidHigh,
				Lo: uuidLow,
			}
			newTimeStamp := timestamppb.Timestamp{
				Seconds: unixSeconds(),
				Nanos:   0,
			}
			newPriceMapOffer := esi.PriceMapOffer{
				Route:    &newRoute,
				OfferId:  &newUuid,
				When:     &newTimeStamp,
				PriceMap: &createdPriceMap,
			}

			err = esi.ProposePriceMapOffer(facilityClient, newPriceMapOffer)
			if err != nil {
				shell.Println(err.Error())
			}
		},
	})

	facilityOffersShellCmd := &ishell.Cmd{
		Name: "offers",
		Help: "manage pending offers",
	}
	shell.AddCmd(facilityOffersShellCmd)
	facilityOffersShellCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "view pending offers",
		Func: func(c *ishell.Context) {
			for k, v := range customerPriceMapOffers {
				shell.Printf("\n%s %s\n%s %s\n%s %s\n%s %s\n",
					boldMsgColorFunc("Customer Public Key:"),
					noteMsgColorFunc(v.Route.GetCustomerKey()),
					boldMsgColorFunc("Producer Public Key:"),
					noteMsgColorFunc(v.Route.GetProducerKey()),
					boldMsgColorFunc("UUID:"),
					k,
					boldMsgColorFunc("Price Map:"),
					v.PriceMap)
			}
			shell.Println()
		},
	})
	facilityOffersShellCmd.AddCmd(&ishell.Cmd{
		Name: "feedback",
		Help: "get feedback on a price map offer",
		Func: func(c *ishell.Context) {
		},
	})
	facilityOffersShellCmd.AddCmd(&ishell.Cmd{
		Name: "evaluate",
		Help: "evaluate an offer and give feedback",
		Func: func(c *ishell.Context) {
			shell.Print("Offer UUID: ")
			uuid, err := uuid2.FromString(c.ReadLine())
			if err != nil {
				shell.Println(err.Error())
				return
			}

			if customerPriceMapOffers[uuid] == nil {
				shell.Printf("no offer with the uuid: '%s'\n", uuid)
				return
			}

			choice := c.MultiChoice([]string{
				"YES",
				"NO",
			}, fmt.Sprintf("Do you accept this offer?\n\n%s\n", customerPriceMapOffers[uuid]))

			newFeedback := esi.PriceMapOfferFeedback{
				Route:   customerPriceMapOffers[uuid].Route,
				OfferId: customerPriceMapOffers[uuid].OfferId,
			}

			if choice == 0 {
				shell.Println("\nOffer has been accepted.\n")
				newFeedback.ObligationStatus = esi.PriceMapOfferFeedback_SATISFIED
				err := esi.ProvidePriceMapOfferFeedback(facilityClient, newFeedback)
				if err != nil {
					shell.Println(err.Error())
				}
			} else if choice == 1 {
				// TODO: counteroffer
			}
		},
	})

	shell.Run()
}

func newPriceMap(shell *ishell.Shell, c *ishell.Context) (esi.PriceMap, error) {
	// Create newPowerComponents.
	shell.Print("Real Power: ")
	realPowerString := c.ReadLine()
	realPower, err := strconv.Atoi(realPowerString)
	if err != nil {
		return esi.PriceMap{}, err
	}
	shell.Print("Reactive Power: ")
	reactivePowerString := c.ReadLine()
	reactivePower, err := strconv.Atoi(reactivePowerString)
	if err != nil {
		return esi.PriceMap{}, err
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
		return esi.PriceMap{}, err
	}
	shell.Print("Expected Duration Nanos: ")
	durationNanosString := c.ReadLine()
	durationNanos, err := strconv.Atoi(durationNanosString)
	if err != nil {
		return esi.PriceMap{}, err
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
		return esi.PriceMap{}, err
	}
	shell.Print("Expected Minimum Duration Nanos: ")
	expectedMinNanosString := c.ReadLine()
	expectedMinNanos, err := strconv.Atoi(expectedMinNanosString)
	if err != nil {
		return esi.PriceMap{}, err
	}
	shell.Print("Expected Maximum Duration Seconds: ")
	expectedMaxSecondsString := c.ReadLine()
	expectedMaxSeconds, err := strconv.Atoi(expectedMaxSecondsString)
	if err != nil {
		return esi.PriceMap{}, err
	}
	shell.Print("Expected Maximum Duration Nanos: ")
	expectedMaxNanosString := c.ReadLine()
	expectedMaxNanos, err := strconv.Atoi(expectedMaxNanosString)
	if err != nil {
		return esi.PriceMap{}, err
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
		return esi.PriceMap{}, err
	}
	shell.Print("Nanos: ")
	nanosString := c.ReadLine()
	nanos, err := strconv.Atoi(nanosString)
	if err != nil {
		return esi.PriceMap{}, err
	}

	// Create a new Price Map.
	//
	// Note that money currency is ignored in this demo.
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

	return newPriceMap, nil
}
