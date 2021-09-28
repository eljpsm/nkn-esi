package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
)

// SignupRegistry discovers facility information, and then sends back a SendKnownDerFacility for each known facility.
func SignupRegistry(client *nkn.MultiClient, registryPublicKey string, info *DerFacilityExchangeInfo) error {
	data, err := proto.Marshal(&RegistryMessage{Chunk: &RegistryMessage_SignupRegistry{SignupRegistry: info}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(registryPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendKnownDerFacility sends facility info.
func SendKnownDerFacility(client *nkn.MultiClient, facilityPublicKey string, info *DerFacilityExchangeInfo) error {
	data, err := proto.Marshal(&FacilityMessage{Chunk: &FacilityMessage_SendKnownDerFacility{SendKnownDerFacility: info}})
	if err != nil {
		return err
	}

	_, err = client.Send(nkn.NewStringArray(facilityPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}

// QueryDerFacilities returns a list of exchanges based on a given location.
func QueryDerFacilities(client *nkn.MultiClient, registryPublicKey string, request *DerFacilityExchangeRequest) error {
	// Encode the given info.
	data, err := proto.Marshal(&RegistryMessage{Chunk: &RegistryMessage_QueryDerFacilities{QueryDerFacilities: request}})
	if err != nil {
		return err
	}

	// Send the information to the Registry.
	_, err = client.Send(nkn.NewStringArray(registryPublicKey), data, nil)
	if err != nil {
		return err
	}

	return nil
}
