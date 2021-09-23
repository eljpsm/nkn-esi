package esi

import (
	"github.com/golang/protobuf/ptypes/empty"
)

// DiscoverRegistry discovers and sends Facility information to a Registry.
func DiscoverRegistry(registryPublicKey string, info DerFacilityExchangeInfo) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// GetPublicKey returns the other Facility's public key.
func GetPublicKey(empty2 empty.Empty) (string, error) {
	return "TODO", nil
}

// GetDerFacilityRegistrationForm returns the registration for a Facility to use.
func GetDerFacilityRegistrationForm(request DerFacilityRegistrationFormRequest) (DerFacilityRegistrationForm, error) {
	return DerFacilityRegistrationForm{}, nil
}

// SubmitDerFacilityRegistrationForm submits a registration form for a Facility.
// When called, the data will be validated, and any problems will be expressed via standard error details.
// When received, the receiving Facility will return with the function CompleteDerFacilityRegistration.
func SubmitDerFacilityRegistrationForm(data DerFacilityRegistrationFormData) (DerFacilityRegistrationFormDataReceipt, error) {
	return DerFacilityRegistrationFormDataReceipt{}, nil
}

// CompleteDerFacilityRegistration completes the Facility registration process.
func CompleteDerFacilityRegistration(registration DerFacilityRegistration) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvideDerCharacteristics publishes DER characteristics for Facilities.
func ProvideDerCharacteristics(characteristics DerCharacteristics) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvidePriceMaps publishes DER price maps for Facilities.
func ProvidePriceMaps(characteristics PriceMapCharacteristics) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProvideSupportedDerPrograms publishes the supported program types.
func ProvideSupportedDerPrograms(set DerProgramSet) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ProposePriceMapOffer propose a price map offer for the service to accept, reject, or propose a counter offer.
// The exchange will invoke this method to make a price map offer to the Facility. THe Facility must respond with either
// an acceptance/rejection of the offer or a counter offer in the form of a different price map proposal.
func ProposePriceMapOffer(request PriceMapOfferStatusRequest) (PriceMapOfferStatus, error) {
	return PriceMapOfferStatus{}, nil
}

// GetPriceMapOfferFeedback returns the status of a price map offer.
func GetPriceMapOfferFeedback(feedback PriceMapOfferFeedback) (PriceMapOfferFeedbackResponse, error) {
	return PriceMapOfferFeedbackResponse{}, nil
}

// ProvidePriceMapOfferStatus provides the status of a price map offer.
func ProvidePriceMapOfferStatus(status PriceMapOfferStatus) (PriceMapOfferStatusResponse, error) {
	return PriceMapOfferStatusResponse{}, nil
}

// ProvidePriceMapOfferFeedback provides feedback on a price map offer, after the offer event is over.
func ProvidePriceMapOfferFeedback(feedback PriceMapOfferFeedback) (PriceMapOfferFeedbackResponse, error) {
	return PriceMapOfferFeedbackResponse{}, nil
}

// ListPrices returns the list of price datum over a time range.
func ListPrices(request DatumRequest) (PriceDatum, error) {
	return PriceDatum{}, nil
}

// ProvidePrices provides pricing data to the Facility.
func ProvidePrices(datum PriceDatum) (empty.Empty, error) {
	return empty.Empty{}, nil
}

// ListPowerProfile returns a list of power profile datum over a time range.
func ListPowerProfile(datum PriceDatum) (PowerProfileDatum, error) {
	return PowerProfileDatum{}, nil
}

// GetPowerParameters gets the power parameters currently used by the services.
func GetPowerParameters(route DerRoute) (PowerParameters, error) {
	return PowerParameters{}, nil
}

// SetPowerParameters sets the power parameters, and then returns the power parameters active after the request.
func SetPowerParameters(parameters PowerParameters) (PowerParameters, error) {
	return PowerParameters{}, nil
}

// GetPriceParameters returns the price parameters currently used by the service.
func GetPriceParameters(route DerRoute) (PriceParameters, error) {
	return PriceParameters{}, nil
}
