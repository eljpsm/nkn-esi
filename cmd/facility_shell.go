package cmd

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	"os"
	"strings"
)

// facilityLoop is the main shell of a Facility.
func facilityShell() error {
	// Create two channels, one for incoming messages and another for outgoing inputs.
	messages := make(chan string)
	inputs := make(chan string)
	go facilityMessageReceiver(messages)
	go facilityInputReceiver(inputs)

	for {
		select {

		case message, ok := <-messages:
			// Evaluate incoming messages.
			if !ok {
				break
			}
			// Print any new message received from the receiver.
			fmt.Println(message)

		case input, ok := <-inputs:
			if !ok {
				break
			}

			err := facilityExecutor(input)
			// If the execution results in an error, alert the user.
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

// facilityMessageReceiver receives and returns any incoming facility messages.
func facilityMessageReceiver(messagesCh chan string) {
	var formKey float64 // a simple number to increment form number.

	message := &esi.FacilityMessage{}

	// TODO: User created? Pass in as argument?
	// An example English language form.
	enForm := esi.Form{
		LanguageCode: "en",
		Key:          fmt.Sprintf("%f", formKey),
		Settings:     nil,
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
			fmt.Println(err.Error())
			continue
		}

		// Case documentation located at api/esi/deer_facility_service.go.
		switch x := message.Chunk.(type) {
		case *esi.FacilityMessage_SendKnownDerFacility:
			// TODO: Does this work as expected for Facility to Facility?
			messagesCh <- fmt.Sprintf("Received Facility from %s - %s", noteMsgColorFunc(msg.Src), infoMsgColorFunc(x.SendKnownDerFacility.GetFacilityPublicKey()))
			knownFacilities[x.SendKnownDerFacility.FacilityPublicKey] = x.SendKnownDerFacility
			messagesCh <- fmt.Sprintf("Saved Facility %s", infoMsgColorFunc(x.SendKnownDerFacility.GetFacilityPublicKey()))

		case *esi.FacilityMessage_GetDerFacilityRegistrationForm:
			// TODO: User created? Pass in as argument?
			messagesCh <- fmt.Sprintf("Received registration from request from %s", noteMsgColorFunc(msg.Src))

			// Fill the registration form with Customer key.
			registrationForm.CustomerFacilityPublicKey = msg.Src

			// Send the registration form.
			err = esi.SendDerFacilityRegistrationForm(facilityClient, registrationForm)
			if err != nil {
				fmt.Println(err.Error())
			}

			formKey += 1 // increment form key
			messagesCh <- fmt.Sprintf("Sent registration form to %s", noteMsgColorFunc(msg.Src))

		case *esi.FacilityMessage_SendDerFacilityRegistrationForm:
			// TODO: User fills in? Automatic? Currently automatic submit.
			messagesCh <- fmt.Sprintf("Received registration form from %s", noteMsgColorFunc(msg.Src))

			// TODO: Fill in the form.
			data := esi.DerFacilityRegistrationFormData{
				CustomerFacilityPublicKey: msg.Src,
			}

			esi.SubmitDerFacilityRegistrationForm(facilityClient, data)

			messagesCh <- fmt.Sprintf("Submitted registration form to %s", noteMsgColorFunc(msg.Src))

		case *esi.FacilityMessage_SubmitDerFacilityRegistrationForm:
			// TODO: Fill in registration form.
			messagesCh <- fmt.Sprintf("Received registration form data from %s", noteMsgColorFunc(msg.Src))

			route := esi.DerRoute{
				BuyKey: facilityInfo.GetFacilityPublicKey(),
				SellKey: msg.Src,
			}
			registration := esi.DerFacilityRegistration{
				Route: &route,
			}

			esi.CompleteDerFacilityRegistration(facilityClient, registration)

			messagesCh <- fmt.Sprintf("Submitted completed registration to %s", noteMsgColorFunc(msg.Src))
			messagesCh <- successMsgColor.Sprintf("Permission granted to %s", noteMsgColorFunc(msg.Src))

		case *esi.FacilityMessage_CompleteDerFacilityRegistration:
			messagesCh <- fmt.Sprintf("Completed registration from %s", noteMsgColorFunc(msg.Src))
			messagesCh <- successMsgColor.Sprintf("Granted permission to %s", noteMsgColorFunc(msg.Src))
		}
	}
}

// facilityInputReceiver receives and returns any facility inputs.
func facilityInputReceiver(inputCh chan string) {
	for {
		input := prompt.Input("> ", facilityCompleter, prompt.OptionPrefixTextColor(prompt.Green))
		inputCh <- input
	}
}

// facilityCompleter is the completer for a Facility.
func facilityCompleter(d prompt.Document) []prompt.Suggest {
	// Useful prompts that the user can use in the shell.
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit out of Facility instance"},
		{Text: "list", Description: "List known facilities"},
		{Text: "signup", Description: "Signup and send Facility info to Registry"},
		{Text: "query", Description: "Query a Registry for Facilities by location"},
		{Text: "request", Description: "Request a registration form from a Facility"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// facilityExecutor is the function which executes user input.
func facilityExecutor(input string) error {
	var err error
	fields := strings.Fields(input)

	if len(fields) == 0 {
		return nil
	}

	// Evaluate the first string.
	switch fields[0] {
	default:
		return errors.New(fmt.Sprintf("unknown command: %s", input))

	case "exit":
		if len(fields) > 1 {
			break
		}

		// Exit out of the program.
		os.Exit(0)

	case "list":
		if len(fields) > 1 {
			break
		}

		for _, v := range knownFacilities {
			fmt.Printf("Name: %s\nCountry: %s\nRegion: %s\nPublic Key: %s\n", v.GetName(), v.Location.GetCountry(), v.Location.GetRegion(), v.GetFacilityPublicKey())
		}

	case "signup":
		for _, v := range fields[1:] {
			// Sign up to a registry.
			err = esi.SignupRegistry(facilityClient, v, facilityInfo)
			if err != nil {
				return err
			}
		}

	case "query":
		// Query a registry by details.
		if len(fields) != 2 {
			break
		}

		myLocation := esi.Location{
			Country: "DC",
		}
		exRequest := esi.DerFacilityExchangeRequest{Location: &myLocation}
		err = esi.QueryDerFacilities(facilityClient, fields[1], exRequest)
		if err != nil {
			return err
		}

	case "request":
		if len(fields) != 3 {
			break
		}

		newRequest := esi.DerFacilityRegistrationFormRequest{FacilityPublicKey: fields[1], LanguageCode: fields[2]}
		err = esi.GetDerFacilityRegistrationForm(facilityClient, newRequest)
		if err != nil {
			return err
		}
	}

	return nil
}
