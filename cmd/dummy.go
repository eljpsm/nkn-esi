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

package cmd

import "github.com/elijahjpassmore/nkn-esi/api/esi"

var (
	// dummyDerRegistryInfo is a dummy DerRegistryInfo.
	dummyDerRegistryInfo = esi.DerRegistryInfo{
		Name:              "New Registry",
		RegistryPublicKey: "",
	}
	// dummyLatLng is a dummy LatLng.
	dummyLatLng = esi.LatLng{
		Latitude:  -36.86397,
		Longitude: 174.72052,
	}
	// dummyLocation is a dummy Location.
	dummyLocation = esi.Location{
		Country:       "New Zealand",
		Region:        "Auckland",
		TimeZone:      "NZT",
		StateProvince: "Auckland",
		PostalCode:    "1022",
		Locality:      "Auckland",
		Sublocality:   "Western Springs",
		StreetAddress: []string{
			"Motions Road",
		},
		Latlng: &dummyLatLng,
	}
	// dummyDerFacilityExchangeInfo is a dummy DerFacilityExchangeInfo.
	dummyDerFacilityExchangeInfo = esi.DerFacilityExchangeInfo{
		Name:              "New Facility",
		FacilityPublicKey: "",
		Location:          &dummyLocation,
	}
)
