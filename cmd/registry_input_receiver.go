package cmd

import (
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
)

// registryInputReceiver receives and returns any registry inputs.
func registryInputReceiver() error {
	message := &esi.RegistryMessage{}

	// Facilities currently stored in memory.
	facilities := make(map[string]*esi.DerFacilityExchangeInfo)

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
			if _, ok := facilities[x.SignupRegistry.FacilityPublicKey]; !ok {
				infoMsgColor.Printf("Saved Facility public key(s) to known Facilities\n")

				facilities[x.SignupRegistry.FacilityPublicKey] = x.SignupRegistry

				for _, v := range facilities {
					esi.SendKnownDerFacility(registryClient, msg.Src, *v)
				}
			}

		case *esi.RegistryMessage_QueryDerFacilities:
			for _, v := range facilities {
				if v.Location.Country == "New Zealand" {
					esi.SendKnownDerFacility(registryClient, msg.Src, *v)
				}
			}
		}
	}
}
