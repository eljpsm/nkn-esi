package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
)

// GetDerFacilityRegistrationForm send a message to get a registration form from request.GetFacilityPublicKey().
func GetDerFacilityRegistrationForm(client *nkn.MultiClient, request DerFacilityRegistrationFormRequest) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_GetDerFacilityRegistrationForm{GetDerFacilityRegistrationForm: &request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.GetPublicKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendDerFacilityRegistrationForm send a message to DerFacilityRegistrationForm to registrationForm.GetCustomerFacilityPublicKey().
func SendDerFacilityRegistrationForm(client *nkn.MultiClient, registrationForm DerFacilityRegistrationForm) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_SendDerFacilityRegistrationForm{SendDerFacilityRegistrationForm: &registrationForm}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registrationForm.GetCustomerFacilityPublicKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SubmitDerFacilityRegistrationForm sends a signed DerFacilityRegistrationFormData to formData.Route.GetSellKey().
func SubmitDerFacilityRegistrationForm(client *nkn.MultiClient, formData DerFacilityRegistrationFormData) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_SubmitDerFacilityRegistrationForm{SubmitDerFacilityRegistrationForm: &formData}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(formData.Route.GetProducerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// CompleteDerFacilityRegistration sends a message to registration.Route.GetBuyKey(), informing them of the registration status.
func CompleteDerFacilityRegistration(client *nkn.MultiClient, registration DerFacilityRegistration) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_CompleteDerFacilityRegistration{CompleteDerFacilityRegistration: &registration}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registration.Route.GetCustomerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetResourceCharacteristics(client *nkn.MultiClient) error {
	return nil
}

// ProposePriceMapOffer propose a price map offer for the service to accept, reject, or propose a counter offer.
// The exchange will invoke this method to make a price map offer to the Facility. The Facility must respond with either
// an acceptance/rejection of the offer or a counter offer in the form of a different price map proposal.
func ProposePriceMapOffer(client *nkn.MultiClient, offer PriceMapOffer) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_ProposePriceMapOffer{ProposePriceMapOffer: &offer}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(offer.Route.GetCustomerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetPriceMapOfferFeedback returns the status of a price map offer.
func GetPriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_GetPriceMapOfferFeedback{GetPriceMapOfferFeedback: &feedback}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(feedback.Route.GetProducerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProvidePriceMapOfferFeedback provides feedback on a price map offer, after the offer event is over.
func ProvidePriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_ProvidePriceMapOfferFeedback{ProvidePriceMapOfferFeedback: &feedback}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(feedback.Route.GetProducerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProvidePrices provides pricing data to the Facility.
func ProvidePrices(client *nkn.MultiClient, datum PriceDatum) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_ProvidePrices{ProvidePrices: &datum}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(datum.Route.GetProducerKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ListPowerProfile returns a list of power profile datum over a time range.
func ListPowerProfile(client *nkn.MultiClient, datum DatumRequest) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_ListPowerProfile{ListPowerProfile: &datum}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(datum.Route.GetProducerKey()), data, nil)
	if err != nil {
		return err
	}

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
