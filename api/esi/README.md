# ESI API

This directory contains the Protobuf based API definition for the concept of Energy Services Interface (**ESI**).

The main entry is through the methods localed within `der_facility_service.go` and `der_facility_registry_service.go`.
There you can find details on sending data. You can then find information on the data received in `der_handler.proto`.

## Previous Work

A previous application of this ESI which leverages the same protobuf structures can be found at
[SolarNetwork](https://github.com/SolarNetwork/der-challenge-prototype), developed for the
[SEPA DER Challenge](https://sepapower.org/plug-and-play-der-challenge/).
