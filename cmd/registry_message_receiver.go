package cmd

import (
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
)

// registryMessageReceiver receives and returns any incoming registry messages.
func registryMessageReceiver() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	<-registryClient.OnConnect.C
	log.WithFields(log.Fields{
		"publicKey": registryInfo.GetRegistryPublicKey(),
		"name": registryInfo.GetName(),
	}).Info("Connection opened")

	message := &esi.RegistryMessage{}

	for {
		msg := <-registryClient.OnMessage.C

		log.WithFields(log.Fields{
			"publicKey": msg.Src,
		}).Info("Message received")

		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// Case documentation located at api/esi/der_facility_registry_service.go.
		switch x := message.Chunk.(type) {
		case *esi.RegistryMessage_SignupRegistry:
			if _, ok := knownFacilities[x.SignupRegistry.FacilityPublicKey]; !ok {
				log.WithFields(log.Fields{
					"publicKey": msg.Src,
				}).Info("Saved facility to known facilities")

				for _, v := range knownFacilities {
					err = esi.SendKnownDerFacility(registryClient, msg.Src, *v)
					if err != nil {
						log.Error(err.Error())
					}
				}

				knownFacilities[x.SignupRegistry.FacilityPublicKey] = x.SignupRegistry
			}

		case *esi.RegistryMessage_QueryDerFacilities:
			// TODO: Look at more than just country.
			log.WithFields(log.Fields{
				"publicKey": msg.Src,
			}).Info("Query for facility")

			for _, v := range knownFacilities {
				if v.Location.Country == x.QueryDerFacilities.Location.GetCountry() {

					// If the facility querying the registry also fits the criteria, ignore it.
					if v.FacilityPublicKey == msg.Src {
						continue
					}

					err = esi.SendKnownDerFacility(registryClient, msg.Src, *v)
					if err != nil {
						log.Error(err.Error())
					}
				}
			}
		}
	}
}
