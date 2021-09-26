package cmd

import (
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
)

// registryMessageReceiver receives and returns any incoming registry messages.
func registryMessageReceiver() {
	message := &esi.RegistryMessage{}

	for {
		msg := <-registryClient.OnMessage.C
		fmt.Printf("Message received from %s\n", noteMsgColorFunc(msg.Src))
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Case documentation located at api/esi/der_facility_registry_service.go.
		switch x := message.Chunk.(type) {
		case *esi.RegistryMessage_SignupRegistry:
			if _, ok := knownFacilities[x.SignupRegistry.FacilityPublicKey]; !ok {
				infoMsgColor.Printf("Saved Facility public key(s) to known Facilities\n")

				for _, v := range knownFacilities {
					esi.SendKnownDerFacility(registryClient, msg.Src, *v)
				}

				knownFacilities[x.SignupRegistry.FacilityPublicKey] = x.SignupRegistry
			}

		case *esi.RegistryMessage_QueryDerFacilities:
			for _, v := range knownFacilities {
				if v.Location.Country == "New Zealand" {

					// If the facility querying the registry also fits the criteria, ignore it.
					if v.FacilityPublicKey == msg.Src {
						continue
					}

					esi.SendKnownDerFacility(registryClient, msg.Src, *v)
				}
			}
		}
	}
}
