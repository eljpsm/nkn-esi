package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nknorg/nkn-sdk-go"
)

// DiscoverRegistry discovers and sends Facility information to a Registry.
func DiscoverRegistry(client *nkn.MultiClient, registryPublicKey string, info DerFacilityExchangeInfo) (empty.Empty, error) {
	// Encode the given info.
	data, err := proto.Marshal(&RegistryMessage{Chunk: &RegistryMessage_Info{Info: &info}})
	if err != nil {
		return empty.Empty{}, err
	}

	// Send the information to the Registry.
	_, err = client.Send(nkn.NewStringArray(registryPublicKey), data, nil)
	if err != nil {
		return empty.Empty{}, err
	}

	return empty.Empty{}, nil
}

// GetDerFacilityRegistrationForm returns the registration for a Facility to use.
func GetDerFacilityRegistrationForm(client *nkn.MultiClient, request DerFacilityRegistrationFormRequest) (DerFacilityRegistrationForm, error) {
	return DerFacilityRegistrationForm{}, nil
}

// SubmitDerFacilityRegistrationForm submits a registration form for a Facility.
// When called, the data will be validated, and any problems will be expressed via standard error details.
// When received, the receiving Facility will return with the function CompleteDerFacilityRegistration.
func SubmitDerFacilityRegistrationForm(client *nkn.MultiClient, data DerFacilityRegistrationFormData) (DerFacilityRegistrationFormDataReceipt, error) {
	return DerFacilityRegistrationFormDataReceipt{}, nil
}

// CompleteDerFacilityRegistration completes the Facility registration process.
func CompleteDerFacilityRegistration(client *nkn.MultiClient, registration DerFacilityRegistration) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvideDerCharacteristics publishes DER characteristics for Facilities.
func ProvideDerCharacteristics(client *nkn.MultiClient, characteristics DerCharacteristics) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvidePriceMaps publishes DER price maps for Facilities.
func ProvidePriceMaps(client *nkn.MultiClient, characteristics PriceMapCharacteristics) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvideSupportedDerPrograms publishes the supported program types.
func ProvideSupportedDerPrograms(client *nkn.MultiClient, set DerProgramSet) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProposePriceMapOffer propose a price map offer for the service to accept, reject, or propose a counter offer.
// The exchange will invoke this method to make a price map offer to the Facility. THe Facility must respond with either
// an acceptance/rejection of the offer or a counter offer in the form of a different price map proposal.
func ProposePriceMapOffer(client *nkn.MultiClient, request PriceMapOfferStatusRequest) (PriceMapOfferStatus, error) {
	return PriceMapOfferStatus{}, nil
}

// GetPriceMapOfferFeedback returns the status of a price map offer.
func GetPriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) (PriceMapOfferFeedbackResponse, error) {
	return PriceMapOfferFeedbackResponse{}, nil
}

// ProvidePriceMapOfferStatus provides the status of a price map offer.
func ProvidePriceMapOfferStatus(client *nkn.MultiClient, status PriceMapOfferStatus) (PriceMapOfferStatusResponse, error) {
	return PriceMapOfferStatusResponse{}, nil
}

// ProvidePriceMapOfferFeedback provides feedback on a price map offer, after the offer event is over.
func ProvidePriceMapOfferFeedback(client *nkn.MultiClient, feedback PriceMapOfferFeedback) (PriceMapOfferFeedbackResponse, error) {
	return PriceMapOfferFeedbackResponse{}, nil
}

// ListPrices returns the list of price datum over a time range.
func ListPrices(client *nkn.MultiClient, request DatumRequest) (PriceDatum, error) {
	return PriceDatum{}, nil
}

// ProvidePrices provides pricing data to the Facility.
func ProvidePrices(client *nkn.MultiClient, datum PriceDatum) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ListPowerProfile returns a list of power profile datum over a time range.
func ListPowerProfile(client *nkn.MultiClient, datum PriceDatum) (PowerProfileDatum, error) {
	return PowerProfileDatum{}, nil
}

// GetPowerParameters gets the power parameters currently used by the services.
func GetPowerParameters(client *nkn.MultiClient, route DerRoute) (PowerParameters, error) {
	return PowerParameters{}, nil
}

// SetPowerParameters sets the power parameters, and then returns the power parameters active after the request.
func SetPowerParameters(client *nkn.MultiClient, parameters PowerParameters) (PowerParameters, error) {
	return PowerParameters{}, nil
}

// GetPriceParameters returns the price parameters currently used by the service.
func GetPriceParameters(client *nkn.MultiClient, route DerRoute) (PriceParameters, error) {
	return PriceParameters{}, nil
}
