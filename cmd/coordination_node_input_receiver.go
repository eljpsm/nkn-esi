package cmd

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

const (
	// defaultLanguage is the default language to look for in a registration form.
	defaultLanguage = "en"
	// defaultCountry is the default country to look for in a query.
	defaultCountry = "DC"

	// defaultRealPower is the default value used for real power when creating a price map.
	defaultRealPower = "10"
	// defaultReactivePower is the default value used for reactive power when creating a price map.
	defaultReactivePower = "10"
	// defaultUnits  is the default value used for units when creating a price map.
	defaultUnits = "100"
	// defaultSeconds is the default value used for time measurement when creating a price map.
	defaultSeconds = 1
)

// coordinationNodeInputReceiver receives and returns any facility inputs.
func coordinationNodeInputReceiver() {
	shell := ishell.New()
	<-coordinationNodeClient.OnConnect.C
	shell.Printf("Connection opened on coordination node '%s'\n", infoMsgColorFunc(coordinationNodeInfo.GetName()))

	coordinationNodeInfoShellCmd := &ishell.Cmd{
		Name: "info",
		Help: "print coordination node information",
	}
	shell.AddCmd(coordinationNodeInfoShellCmd)
	coordinationNodeInfoShellCmd.AddCmd(&ishell.Cmd{
		Name: "public",
		Help: "print local public key of coordination node",
		Func: func(c *ishell.Context) {
			shell.Printf("%s\n", infoMsgColorFunc(coordinationNodeInfo.GetPublicKey()))
		},
	})

	coordinationNodeRegistryShellCmd := &ishell.Cmd{
		Name: "registry",
		Help: "manage registry functionality",
	}
	shell.AddCmd(coordinationNodeRegistryShellCmd)
	coordinationNodeRegistryShellCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "print known coordination nodes received from registry",
		Func: func(c *ishell.Context) {
			for _, facility := range knownCoordinationNodes {
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
	coordinationNodeRegistryShellCmd.AddCmd(&ishell.Cmd{
		Name: "signup",
		Help: "sign up to a registry",
		Func: func(c *ishell.Context) {
			c.Print("Registry Public Key: ")

			err := esi.SignupRegistry(coordinationNodeClient, c.ReadLine(), &coordinationNodeInfo)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})
	coordinationNodeRegistryShellCmd.AddCmd(&ishell.Cmd{
		Name: "query",
		Help: "query registry for coordination nodes by location",
		Func: func(c *ishell.Context) {
			c.Print("Registry Public Key: ")
			registryPublicKey := c.ReadLine()
			c.Printf("Country [%s]: ", defaultCountry)
			country := c.ReadLine()
			if country == "" {
				country = defaultCountry
			}

			// You can query based upon any setting that DerFacilityExchangeRequest takes.
			//
			// In this demo, you can select a COUNTRY to query based upon.
			newLocation := esi.Location{
				Country: country,
			}

			request := esi.DerFacilityExchangeRequest{Location: &newLocation}
			err := esi.QueryDerFacilities(coordinationNodeClient, registryPublicKey, &request)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "peers",
		Help: "show any registered facilities or exchanges",
		Func: func(c *ishell.Context) {
			if registeredExchange != "" {
				// Print the exchange.
				shell.Printf("\n%s\n", boldMsgColorFunc("EXCHANGE"))
				shell.Printf("%s %s\n",
					boldMsgColorFunc("Public Key:"),
					noteMsgColorFunc(registeredExchange))
			}
			if len(registeredFacilities) > 0 {
				// Print the facilities.
				shell.Printf("\n%s\n", boldMsgColorFunc("FACILITIES"))
				for k := range registeredFacilities {
					shell.Printf("%s %s\n",
						boldMsgColorFunc("Public Key:"),
						noteMsgColorFunc(k))

					if facilityCharacteristics[k] != nil {
						shell.Println(proto.MarshalTextString(facilityCharacteristics[k]))
					}
					if facilityPriceMaps[k] != nil {
						shell.Println(proto.MarshalTextString(facilityPriceMaps[k]))
					}
				}
				shell.Println()
			} else {
				shell.Println()
			}
		},
	})

	coordinationNodeFacilityShellCmd := &ishell.Cmd{
		Name: "facility",
		Help: "manage coordination node facility functionality",
	}
	shell.AddCmd(coordinationNodeFacilityShellCmd)
	coordinationNodeFacilityShellCmd.AddCmd(&ishell.Cmd{
		Name: "request",
		Help: "request registration form from a coordination node behaving as an exchange",
		Func: func(c *ishell.Context) {
			if registeredExchange != "" {
				shell.Println("customer has already been set")
				return
			}
			c.Print("Public Key: ")
			exchangePublicKey := c.ReadLine()
			c.Printf("Language Code [%s]: ", defaultLanguage)
			languageCode := c.ReadLine()
			if languageCode == "" {
				languageCode = defaultLanguage
			}

			// When creating a request, you can specify a language code.
			//
			// In this demo, the only language code that is used (and sent) is "en" for English.
			request := esi.DerFacilityRegistrationFormRequest{
				PublicKey:    exchangePublicKey,
				LanguageCode: languageCode,
			}

			err := esi.GetDerFacilityRegistrationForm(coordinationNodeClient, &request)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})
	coordinationNodeFacilityShellCmd.AddCmd(&ishell.Cmd{
		Name: "forms",
		Help: "print forms to be signed",
		Func: func(c *ishell.Context) {
			if registeredExchange != "" {
				shell.Println("customer has already been set")
				return
			}
			for _, v := range receivedRegistrationForms {
				shell.Printf("%s %s",
					boldMsgColorFunc("Exchange Public Key:"),
					noteMsgColorFunc(v.Route.GetExchangeKey()))
			}
			shell.Println()
		},
	})
	coordinationNodeFacilityShellCmd.AddCmd(&ishell.Cmd{
		Name: "register",
		Help: "fill in a received registration form",
		Func: func(c *ishell.Context) {
			if registeredExchange != "" {
				shell.Println("customer has already been set")
				return
			}
			shell.Print("Public Key: ")
			publicKey := c.ReadLine()

			form, present := receivedRegistrationForms[publicKey]

			// TODO: needed? check
			if present {
				shell.Println() // gap from input

				// TODO: nonce
				// Contains the results of key -> response.
				results := make(map[string]string)
				route := esi.DerRoute{
					ExchangeKey: form.Route.GetExchangeKey(),
					FacilityKey: form.Route.GetFacilityKey(),
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
					// situations in which user input could be either optional or unnecessary.
					if setting.GetPlaceholder() != "" {
						if result == "" {
							result = setting.GetPlaceholder()
						}
					}

					results[setting.Key] = result
				}

				// Submit the registration form.
				err := esi.SubmitDerFacilityRegistrationForm(coordinationNodeClient, &registrationFormData)
				if err != nil {
					log.Error(err.Error())
				}

				// Remove form from the map.
				delete(receivedRegistrationForms, form.Route.GetExchangeKey())

				log.WithFields(log.Fields{
					"end": form.Route.GetExchangeKey(),
				}).Info("Sent registration form")

				shell.Printf("\nForm has been submitted to %s\n", registrationFormData.Route.GetExchangeKey())

			} else {
				shell.Printf("no form found with public key '%s`\n", publicKey)
				return
			}
		},
	})

	coordinationNodePriceMapShellCmd := &ishell.Cmd{
		Name: "price-map",
		Help: "manage local price maps",
	}
	shell.AddCmd(coordinationNodePriceMapShellCmd)
	coordinationNodePriceMapShellCmd.AddCmd(&ishell.Cmd{
		Name: "view",
		Help: "print local price map",
		Func: func(c *ishell.Context) {
			fmt.Println(proto.MarshalTextString(&priceMap))
		},
	})
	coordinationNodePriceMapShellCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "create a local price map",
		Func: func(c *ishell.Context) {
			createdPriceMap, err := newPriceMap(shell, c, defaultRealPower, defaultReactivePower, defaultUnits)
			if err != nil {
				shell.Println(err.Error())
				return
			}

			// TODO: fix
			priceMap = *createdPriceMap
		},
	})

	coordinationNodeCharacteristicsShellCmd := &ishell.Cmd{
		Name: "characteristics",
		Help: "manage local characteristics",
	}
	shell.AddCmd(coordinationNodeCharacteristicsShellCmd)
	coordinationNodeCharacteristicsShellCmd.AddCmd(&ishell.Cmd{
		Name: "view",
		Help: "print local characteristics",
		Func: func(c *ishell.Context) {
			fmt.Println(&resourceCharacteristics)
		},
	})
	coordinationNodeCharacteristicsShellCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "create coordination node facility characteristics",
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

	coordinationNodeExchangeShellCmd := &ishell.Cmd{
		Name: "exchange",
		Help: "manage coordination node exchange functionality",
	}
	coordinationNodeExchangeShellCmd.AddCmd(&ishell.Cmd{
		Name: "get",
		Help: "get characteristics and price map of coordination node behaving as a facility",
		Func: func(c *ishell.Context) {
			shell.Print("Public Key: ")
			publicKey := c.ReadLine()

			if !registeredFacilities[publicKey] {
				shell.Printf("no facility with public key: '%s\n'", publicKey)
				return
			}

			newRoute := esi.DerRoute{
				ExchangeKey: coordinationNodeInfo.GetPublicKey(),
				FacilityKey: publicKey,
			}
			newCharacteristicsRequest := esi.DerResourceCharacteristicsRequest{
				Route: &newRoute,
			}
			newPriceMapRequest := esi.DerPriceMapRequest{
				Route: &newRoute,
			}

			err := esi.GetResourceCharacteristics(coordinationNodeClient, &newCharacteristicsRequest)
			if err != nil {
				log.Error(err.Error())
			}
			err = esi.GetPriceMap(coordinationNodeClient, &newPriceMapRequest)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})
	coordinationNodeExchangeShellCmd.AddCmd(&ishell.Cmd{
		Name: "price-maps",
		Help: "print facility price maps",
		Func: func(c *ishell.Context) {
			for k, v := range facilityPriceMaps {
				shell.Printf("\n%s %s\n%s %s\n\n",
					boldMsgColorFunc("Public Key:"),
					noteMsgColorFunc(k),
					boldMsgColorFunc("Price Map:"),
					v)
			}
		},
	})
	coordinationNodeExchangeShellCmd.AddCmd(&ishell.Cmd{
		Name: "characteristics",
		Help: "print facility characteristics",
		Func: func(c *ishell.Context) {
			for k, v := range facilityCharacteristics {
				shell.Printf("\n%s %s\n%s %s\n\n",
					boldMsgColorFunc("Public Key:"),
					noteMsgColorFunc(k),
					boldMsgColorFunc("Characteristics:"),
					v)
			}
		},
	})
	shell.AddCmd(coordinationNodeExchangeShellCmd)
	coordinationNodeExchangeShellCmd.AddCmd(&ishell.Cmd{
		Name: "propose",
		Help: "propose a price map offer to a coordination node behaving as facility",
		Func: func(c *ishell.Context) {
			shell.Print("Public Key: ")
			publicKey := c.ReadLine()
			if !registeredFacilities[publicKey] {
				shell.Printf("no facility with public key: '%s'\n", publicKey)
				return
			}

			createdPriceMap, err := newPriceMap(shell, c, defaultRealPower, defaultReactivePower, defaultUnits)
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newRoute := esi.DerRoute{
				FacilityKey: publicKey,
				ExchangeKey: coordinationNodeInfo.GetPublicKey(),
			}
			uuid, err := newUuid()
			if err != nil {
				shell.Println(err.Error())
				return
			}
			newUuid := esi.Uuid{
				Uuid: uuid,
			}
			newTimeStamp := timestamppb.Timestamp{
				Seconds: unixSeconds(),
				Nanos:   0,
			}
			newPriceMapOffer := esi.PriceMapOffer{
				Route:    &newRoute,
				OfferId:  &newUuid,
				When:     &newTimeStamp,
				PriceMap: createdPriceMap,
			}

			priceMapOffers[uuid] = &newPriceMapOffer
			// Store the status of the offer.
			status := esi.PriceMapOfferStatus{
				Route:   &newRoute,
				OfferId: &newUuid,
				Status:  0, // store unknown status
			}
			priceMapOfferStatus[uuid] = &status

			err = esi.ProposePriceMapOffer(coordinationNodeClient, &newPriceMapOffer)
			if err != nil {
				log.Error(err.Error())
			}
		},
	})

	coordinationNodeOffersShellCmd := &ishell.Cmd{
		Name: "offers",
		Help: "manage pending offers",
	}
	shell.AddCmd(coordinationNodeOffersShellCmd)
	coordinationNodeOffersShellCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "view all offers",
		Func: func(c *ishell.Context) {
			for k, v := range priceMapOffers {
				// You have access to a lot of information.
				//
				// In this example, only key information is provided.
				shell.Printf("\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n",
					boldMsgColorFunc("Exchange Public Key:"),
					noteMsgColorFunc(v.Route.GetExchangeKey()),
					boldMsgColorFunc("Facility Public Key:"),
					noteMsgColorFunc(v.Route.GetFacilityKey()),
					boldMsgColorFunc("UUID:"),
					k,
					boldMsgColorFunc("Price Map:"),
					proto.MarshalTextString(v.PriceMap),
					boldMsgColorFunc("Status:"),
					infoMsgColorFunc(priceMapOfferStatus[v.OfferId.Uuid].Status))
			}
			shell.Println()
		},
	})
	coordinationNodeOffersShellCmd.AddCmd(&ishell.Cmd{
		Name: "evaluate",
		Help: "evaluate an offer and give a response",
		Func: func(c *ishell.Context) {
			shell.Print("Offer UUID: ")
			currentUuid := c.ReadLine()

			if priceMapOffers[currentUuid] == nil {
				shell.Printf("no offer with the uuid: '%s'\n", currentUuid)
				return
			}

			choice := c.MultiChoice([]string{
				"YES",
				"NO",
			}, fmt.Sprintf("Do you accept this offer?\n\n%s\n", proto.MarshalTextString(priceMapOffers[currentUuid])))

			if choice == 0 {
				// Accept the offer.
				shell.Println("\nOffer has been accepted.\n")
				accept := esi.PriceMapOfferResponse_Accept{
					Accept: true,
				}
				response := esi.PriceMapOfferResponse{
					Route:       priceMapOffers[currentUuid].Route,
					OfferId:     priceMapOffers[currentUuid].OfferId,
					AcceptOneof: &accept,
				}
				err := esi.SendPriceMapOfferResponse(coordinationNodeClient, &response)
				if err != nil {
					log.Error(err.Error())
				}

				priceMap = *priceMapOffers[currentUuid].PriceMap
				// Store the status ACCEPTED.
				priceMapOfferStatus[currentUuid].Status = 1

				log.WithFields(log.Fields{
					"src": priceMapOffers[currentUuid].Route.GetExchangeKey(),
				}).Info("Accepted price map")
				log.Info("Updated price map")

			} else if choice == 1 {
				// Create a new counter offer.
				//
				// In reality, this process may be more sophisticated - but for this demo, you will keep sending
				// counter offers until one is accepted.
				shell.Println()
				createdPriceMap, err := newPriceMap(
					shell,
					c,
					strconv.FormatInt(priceMapOffers[currentUuid].PriceMap.PowerComponents.RealPower, 10),
					strconv.FormatInt(priceMapOffers[currentUuid].PriceMap.PowerComponents.ReactivePower, 10),
					strconv.FormatInt(priceMapOffers[currentUuid].PriceMap.Price.ApparentEnergyPrice.Units, 10))
				if err != nil {
					shell.Println(err.Error())
					return
				}
				counterOffer := esi.PriceMapOfferResponse_CounterOffer{
					CounterOffer: createdPriceMap,
				}
				uuid, err := newUuid()
				if err != nil {
					shell.Println(err.Error())
					return
				}
				newUuid := esi.Uuid{
					Uuid: uuid,
				}
				offerResponse := esi.PriceMapOfferResponse{
					Route:         priceMapOffers[currentUuid].Route,
					PreviousOffer: priceMapOffers[currentUuid].OfferId,
					OfferId:       &newUuid,
					AcceptOneof:   &counterOffer,
				}

				err = esi.SendPriceMapOfferResponse(coordinationNodeClient, &offerResponse)
				if err != nil {
					log.Error(err.Error())
				}

				// Store the status REJECTED.
				priceMapOfferStatus[currentUuid].Status = 2
			}
		},
	})

	shell.Run()
}

// newPriceMap creates and returns a new price map.
func newPriceMap(shell *ishell.Shell, c *ishell.Context, optRealPower string, optReactivePower string, optUnits string) (*esi.PriceMap, error) {
	// Create newPowerComponents.
	shell.Printf("Real Power [%s]: ", optRealPower)
	realPowerString := c.ReadLine()
	if realPowerString == "" {
		realPowerString = optRealPower
	}
	realPower, err := strconv.Atoi(realPowerString)
	if err != nil {
		return &esi.PriceMap{}, err
	}

	shell.Printf("Reactive Power [%s]: ", optReactivePower)
	reactivePowerString := c.ReadLine()
	if reactivePowerString == "" {
		reactivePowerString = optReactivePower
	}
	reactivePower, err := strconv.Atoi(reactivePowerString)
	if err != nil {
		return &esi.PriceMap{}, err
	}

	newPowerComponents := esi.PowerComponents{
		RealPower:     int64(realPower),
		ReactivePower: int64(reactivePower),
	}

	// Create newDuration.
	//
	// In this demo, all price maps are immediately completed for simplicity.
	durationSeconds := 0
	durationNanos := 0
	newDuration := duration.Duration{
		Seconds: int64(durationSeconds),
		Nanos:   int32(durationNanos),
	}

	// Create newDurationRange.
	//
	// In this demo, the response time is assumed to be immediate.
	expectedMinSeconds := defaultSeconds
	expectedMinNanos := defaultSeconds
	expectedMaxSeconds := defaultSeconds
	expectedMaxNanos := defaultSeconds
	if err != nil {
		return &esi.PriceMap{}, err
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

	// Currency currently isn't evaluated, so just set it to USD.
	currencyCode := "USD"
	shell.Printf("Units [%s]: ", optUnits)
	unitsString := c.ReadLine()
	if unitsString == "" {
		unitsString = optUnits
	}
	units, err := strconv.Atoi(unitsString)
	if err != nil {
		return &esi.PriceMap{}, err
	}
	// Don't worry about setting nanos.
	nanos := 0

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

	return &newPriceMap, nil
}
