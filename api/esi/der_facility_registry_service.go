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

// der_facility_registry_service.go
//
// The functions contained here can be thought of as the calling functions.
//
// E.g:
//		SignupRegistry(...)
//
// should be read as:
//		Signup my coordination node to registry ...
//
// For information on returning behaviour, consult der_handler.go.

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
