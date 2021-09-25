package cmd

import (
	"fmt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
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

		// Case documentation located at api/esi/der_handler.proto.
		switch x := message.Chunk.(type) {
		case *esi.RegistryMessage_DerFacilityExchangeInfo:
			if _, ok := facilities[x.DerFacilityExchangeInfo.FacilityPublicKey]; !ok {
				infoMsgColor.Printf("Saved Facility public key(s) to known Facilities\n")

				facilities[x.DerFacilityExchangeInfo.FacilityPublicKey] = x.DerFacilityExchangeInfo

				for _, v := range facilities {
					data, err := proto.Marshal(&esi.FacilityMessage{Chunk: &esi.FacilityMessage_DerFacilityExchangeInfo{DerFacilityExchangeInfo: v}})
					if err != nil {
						panic(err)
					}

					_, err = registryClient.Send(nkn.NewStringArray(msg.Src), data, nil)
					if err != nil {
						panic(err)
					}
				}
			}

		case *esi.RegistryMessage_DerFacilityExchangeRequest:
			for _, v := range facilities {
				if v.Location.Country == "New Zealand" {
					data, _ := proto.Marshal(&esi.FacilityMessage{Chunk: &esi.FacilityMessage_DerFacilityExchangeInfo{DerFacilityExchangeInfo: v}})
					fmt.Printf("Send Facility %s to %s\n", infoMsgColorFunc(v.FacilityPublicKey), noteMsgColorFunc(msg.Src))
					_, err = registryClient.Send(nkn.NewStringArray(msg.Src), data, nil)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
