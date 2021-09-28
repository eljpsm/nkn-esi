package cmd

import (
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// registryMessageReceiver receives and returns any incoming registry messages.
func registryMessageReceiver() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	<-registryClient.OnConnect.C
	log.WithFields(log.Fields{
		"publicKey": registryInfo.GetPublicKey(),
		"name":      registryInfo.GetName(),
	}).Info("Connection opened")

	message := &esi.RegistryMessage{}

	for {
		msg := <-registryClient.OnMessage.C

		log.WithFields(log.Fields{
			"src": msg.Src,
		}).Info("Message received")

		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// Case documentation located at api/esi/der_facility_registry_service.go.
		switch x := message.Chunk.(type) {
		case *esi.RegistryMessage_SignupRegistry:
			if _, ok := knownCoordinationNodes[x.SignupRegistry.PublicKey]; !ok {
				log.WithFields(log.Fields{
					"src": msg.Src,
				}).Info("Saved coordination node")

				for _, facility := range knownCoordinationNodes {
					err = esi.SendKnownDerFacility(registryClient, msg.Src, facility)
					if err != nil {
						log.Error(err.Error())
					}
				}

				knownCoordinationNodes[x.SignupRegistry.PublicKey] = x.SignupRegistry
			}

		case *esi.RegistryMessage_QueryDerFacilities:
			// TODO: Look at more than just country.
			log.WithFields(log.Fields{
				"src": msg.Src,
			}).Info("Query for coordination node")

			for _, coordinationNode := range knownCoordinationNodes {
				if strings.ToLower(coordinationNode.Location.GetCountry()) == strings.ToLower(x.QueryDerFacilities.Location.GetCountry()) {

					// If the facility querying the registry also fits the criteria, ignore it.
					if coordinationNode.PublicKey == msg.Src {
						continue
					}

					err = esi.SendKnownDerFacility(registryClient, msg.Src, coordinationNode)
					if err != nil {
						log.Error(err.Error())
					}

					log.WithFields(log.Fields{
						"end": msg.Src,
					}).Info("Sent known coordination node")
				}
			}
		}
	}
}
