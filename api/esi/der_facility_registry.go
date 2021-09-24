package esi

import (
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
)

// ListDerFacilities returns a list of exchanges based on a given location.
func ListDerFacilities(client *nkn.MultiClient, registryPublicKey string, request DerFacilityExchangeRequest) error {
	// Encode the given info.
	data, err := proto.Marshal(&RegistryMessage{Chunk: &RegistryMessage_List{List: &request}})
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
