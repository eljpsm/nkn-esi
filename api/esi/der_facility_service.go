package esi

import "github.com/golang/protobuf/ptypes/empty"

func CompleteDerFacilityRegistration(registration DerFacilityRegistration) (empty.Empty, error) {
	return empty.Empty{}, nil
}

func ProposePriceMapOfficer(request PriceMapOfferStatusRequest) (PriceMapOfferStatus, error) {
	return PriceMapOfferStatus{}, nil
}

func ProvidePriceMapOfferFeedback(feedback PriceMapOfferFeedback) (PriceMapOfferFeedbackResponse, error) {
	return PriceMapOfferFeedbackResponse{}, nil
}

func ProvidePrices(datum PriceDatum) (empty.Empty, error) {
	return empty.Empty{}, nil
}

func ListPowerProfile(datum PriceDatum) (PowerProfileDatum, error) {
	return PowerProfileDatum{}, nil
}

func GetPowerParameters(route DerRoute) (PowerParameters, error) {
	return PowerParameters{}, nil
}

func SetPowerParameters(parameters PowerParameters) (PowerParameters, error) {
	return PowerParameters{}, nil
}

func GetPriceParameters(route DerRoute) (PriceParameters, error) {
	return PriceParameters{}, nil
}
