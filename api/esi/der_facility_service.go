package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
)

// GetDerFacilityRegistrationForm returns the registration for a Facility to use.
func GetDerFacilityRegistrationForm(client *nkn.MultiClient, request DerFacilityRegistrationFormRequest) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_GetDerFacilityRegistrationForm{GetDerFacilityRegistrationForm: &request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.FacilityPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendDerFacilityRegistrationForm sends the registration form to the customer.
func SendDerFacilityRegistrationForm(client *nkn.MultiClient, registrationForm DerFacilityRegistrationForm) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_SendDerFacilityRegistrationForm{SendDerFacilityRegistrationForm: &registrationForm}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registrationForm.CustomerFacilityPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SubmitDerFacilityRegistrationForm submits a registration form for a Facility.
// When called, the data will be validated, and any problems will be expressed via standard error details.
// When received, the receiving Facility will return with the function CompleteDerFacilityRegistration.
func SubmitDerFacilityRegistrationForm(client *nkn.MultiClient, formData DerFacilityRegistrationFormData) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_SubmitDerFacilityRegistrationForm{SubmitDerFacilityRegistrationForm: &formData}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(formData.CustomerFacilityPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// CompleteDerFacilityRegistration completes the Facility registration process.
func CompleteDerFacilityRegistration(client *nkn.MultiClient, registration DerFacilityRegistration) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_CompleteDerFacilityRegistration{CompleteDerFacilityRegistration: &registration}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registration.Route.BuyKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProvideDerCharacteristics publishes DER characteristics for Facilities.
func ProvideDerCharacteristics(client *nkn.MultiClient, characteristics DerCharacteristics) error {
	return nil
}

// ProvidePriceMaps publishes DER price maps for Facilities.
func ProvidePriceMaps(client *nkn.MultiClient, characteristics PriceMapCharacteristics) error {
	return nil
}

// ProvideSupportedDerPrograms publishes the supported program types.
func ProvideSupportedDerPrograms(client *nkn.MultiClient, set DerProgramSet) error {
	return nil
}

// ProposePriceMapOffer propose a price map offer for the service to accept, reject, or propose a counter offer.
// The exchange will invoke this method to make a price map offer to the Facility. THe Facility must respond with either
// an acceptance/rejection of the offer or a counter offer in the form of a different price map proposal.
func ProposePriceMapOffer(client *nkn.MultiClient, request PriceMapOfferStatusRequest) error {
	return nil
}

// GetPriceMapOfferFeedback returns the status of a price map offer.
func GetPriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) error {
	return nil
}

// ProvidePriceMapOfferStatus provides the status of a price map offer.
func ProvidePriceMapOfferStatus(client *nkn.MultiClient, status PriceMapOfferStatus) error {
	return nil
}

// ProvidePriceMapOfferFeedback provides feedback on a price map offer, after the offer event is over.
func ProvidePriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) error {
	return nil
}

// ListPrices returns the list of price datum over a time range.
func ListPrices(client *nkn.MultiClient, request DatumRequest) error {
	return nil
}

// ProvidePrices provides pricing data to the Facility.
func ProvidePrices(client *nkn.MultiClient, datum PriceDatum) error {
	return nil
}

// ListPowerProfile returns a list of power profile datum over a time range.
func ListPowerProfile(client *nkn.MultiClient, datum DatumRequest) error {
	return nil
}

// GetPowerParameters gets the power parameters currently used by the services.
func GetPowerParameters(client *nkn.MultiClient, route DerRoute) error {
	return nil
}

// SetPowerParameters sets the power parameters, and then returns the power parameters active after the request.
func SetPowerParameters(client *nkn.MultiClient, parameters PowerParameters) error {
	return nil
}

// GetPriceParameters returns the price parameters currently used by the service.
func GetPriceParameters(client *nkn.MultiClient, route DerRoute) error {
	return nil
}
