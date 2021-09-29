/*
Copyright © 2021 Ecogy Energy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
)

// coordination_node_service.go
//
// This implements the functionality defined by the facility and exchange ESI API.
//
// The functions contained here can be thought of as the calling functions.
//
// E.g:
//		SubmitDerFacilityRegistrationForm(...)
//
// should be read as:
//		Submit my facility registration form to ...
//
// For information on returning behaviour, consult der_handler.go.

// GetDerFacilityRegistrationForm sends a message to an exchange to receive a registration form.
func GetDerFacilityRegistrationForm(client *nkn.MultiClient, request *DerFacilityRegistrationFormRequest) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_GetDerFacilityRegistrationForm{GetDerFacilityRegistrationForm: request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.GetPublicKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendDerFacilityRegistrationForm sends a facility registration form to a facility.
func SendDerFacilityRegistrationForm(client *nkn.MultiClient, registrationForm *DerFacilityRegistrationForm) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SendDerFacilityRegistrationForm{SendDerFacilityRegistrationForm: registrationForm}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registrationForm.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SubmitDerFacilityRegistrationForm sends a completed facility registration form to an exchange.
func SubmitDerFacilityRegistrationForm(client *nkn.MultiClient, formData *DerFacilityRegistrationFormData) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SubmitDerFacilityRegistrationForm{SubmitDerFacilityRegistrationForm: formData}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(formData.Route.GetExchangeKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// CompleteDerFacilityRegistration sends a notification to a facility of a successful registration.
func CompleteDerFacilityRegistration(client *nkn.MultiClient, registration *DerFacilityRegistration) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_CompleteDerFacilityRegistration{CompleteDerFacilityRegistration: registration}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registration.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetResourceCharacteristics sends a request for facility resource characteristics.
func GetResourceCharacteristics(client *nkn.MultiClient, request *DerResourceCharacteristicsRequest) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_GetResourceCharacteristics{GetResourceCharacteristics: request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendResourceCharacteristics sends resource characteristics to the exchange.
func SendResourceCharacteristics(client *nkn.MultiClient, characteristics *DerCharacteristics) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SendResourceCharacteristics{SendResourceCharacteristics: characteristics}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(characteristics.Route.GetExchangeKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetPriceMap sends a request for the facility price map.
func GetPriceMap(client *nkn.MultiClient, request *DerPriceMapRequest) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_GetPriceMap{GetPriceMap: request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendPriceMap sends the price map to the exchange.
func SendPriceMap(client *nkn.MultiClient, exchangeKey string, priceMap *PriceMap) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SendPriceMap{SendPriceMap: priceMap}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(exchangeKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProposePriceMapOffer proposes a price map offer for the other party to accept, reject, or propose a counter offer.
// The exchange will invoke this method to make a price map offer to the Facility. The Facility must respond with either
// an acceptance/rejection of the offer or a counter offer in the form of a different price map proposal.
//
// This function will optionally switch the node type if provided. This allows systems which combine facility and
// exchange behaviour into one to more easily manage routing.
func ProposePriceMapOffer(client *nkn.MultiClient, offer *PriceMapOffer) error {
	var address string
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_ProposePriceMapOffer{ProposePriceMapOffer: offer}})
	if err != nil {
		return err
	}

	if offer.Node.Type == NodeType_FACILITY {
		address = offer.Route.GetFacilityKey()
	} else {
		address = offer.Route.GetExchangeKey()
	}

	_, err = client.Send(nkn.NewStringArray(address), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendPriceMapOfferResponse sends an offer response to the other party
//
// This function will optionally switch the node type if provided. This allows systems which combine facility and
// exchange behaviour into one to more easily manage routing.
func SendPriceMapOfferResponse(client *nkn.MultiClient, response *PriceMapOfferResponse) error {
	var address string
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SendPriceMapOfferResponse{SendPriceMapOfferResponse: response}})
	if err != nil {
		return err
	}

	if response.Node.Type == NodeType_FACILITY {
		address = response.Route.GetFacilityKey()
	} else {
		address = response.Route.GetExchangeKey()
	}

	_, err = client.Send(nkn.NewStringArray(address), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetPriceMapOfferFeedback sends offer feedback to the exchange to return a feedback response.
func GetPriceMapOfferFeedback(client *nkn.MultiClient, feedback *PriceMapOfferFeedback) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_GetPriceMapOfferFeedback{GetPriceMapOfferFeedback: feedback}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(feedback.Route.GetExchangeKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProvidePriceMapOfferFeedback provides feedback on a price map offer, after the offer event is over.
func ProvidePriceMapOfferFeedback(client *nkn.MultiClient, response *PriceMapOfferFeedbackResponse) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_ProvidePriceMapOfferFeedback{ProvidePriceMapOfferFeedback: response}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(response.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetPowerParameters gets the power parameters currently used by the services.
func GetPowerParameters(client *nkn.MultiClient, request *DerPowerParametersRequest) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_GetPowerParameters{GetPowerParameters: request}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(request.Route.GetExchangeKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetPowerParameters sets the power parameters.
func SetPowerParameters(client *nkn.MultiClient, facilityKey string, parameters *PowerParameters) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_SetPowerParameters{SetPowerParameters: parameters}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(facilityKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// ListPrices sends regular location based price datum to a facility.
func ListPrices(client *nkn.MultiClient, datum *PriceDatum) error {
	data, err := proto.Marshal(&CoordinationNodeMessage{Chunk: &CoordinationNodeMessage_ListPrices{ListPrices: datum}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(datum.Route.GetFacilityKey()), data, nil)
	if err != nil {
		return err
	}

	return nil
}
