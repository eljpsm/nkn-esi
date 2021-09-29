# nkn-esi

NKN-ESI (or **nESI**) is an [NKN](https://nkn.org/) based Energy Services Interface (ESI). An ESI supports a distributed
marketplace for energy services on the electricity grid. It can be used to facilitate services such as load shifting
(e.g. delaying energy consumption to help with peak capacity management) or the timed increased consumption of energy
(e.g. activating devices to consume energy to mitigate a high voltage situation). The purpose of these services is to
allow an aggregator, utility, or distribution system operator to easily and cost effectively maintain a stable and
resilient electricity grid.

By developing NKN-ESI, we hope to show that an ESI which leverages NKN could be a very simple, secure, and resilient
mechanism for building, factory, facility, or distributed energy resource operators to offer and provide energy services
to a range of entities responsible for aspects of grid stability.

## Previous Work

A previous application of this ESI can be found at
[SolarNetwork](https://github.com/SolarNetwork/der-challenge-prototype), developed for the
[SEPA DER Challenge](https://sepapower.org/plug-and-play-der-challenge/).

## Installation

To install nkn-esi, you will need **go**, **make** and **protoc**.

You can find the latest version of Go at [golang.org](https://golang.org/doc/install).

### Ubuntu

```bash
sudo apt-get install golang-goprotobuf-dev
sudo apt-get install build-essential
```

1. Execute `make` in the project root directory.
2. Execute `go build` in the project root directory.

## Demo

### Conceptual

This demo utilizes the ESI API to create two primary agents:

* the **coordination node**
* the **registry**

The coordination node is a functional concatenation of the two ESI concepts of the **facility**
(the DER Facility, or DERF), and the **exchange** (the Interfacing Party with External Responsibility, or IPER).  This
coordination node allows you to use a single instance as both a facility to an exchange, and an exchange to facilities.

The registry is a simple server that allows exchanges to "signup" and save their information to. Then, a facility can
query the registry looking for specific exchange details and easily match. A registry is decentralized, and anyone can
create their own registry and have exchanges sign up.

### Setup

For the purposes of this demo, you will start three separate instances:

1. the registry
2. the facility
3. the exchange

Note that both the facility and exchange are instances of the coordination node.

Create a directory in the project root directory called `configs`. This is ignored by git, so you can place any of your
configuration or secret files in there.

#### Registry

First, you must create a configuration of your registry.

1. Run `./nkn-esi registry init configs/registry`. This will create two files: `registry.json` (your configuration),
   and `registry.secret` (your secret key). Your registry configuration is fine to share with others, but your secret
   key should be kept to yourself. Leave the configuration file as it is, the defaults are fine.
2. Run `./nkn-esi registry start configs/registry.json configs/registry.secret`. This will load your config and secret
   key for this particular instance.

You're done! You will now see a log print out that the registry has started.

#### Facility

You must now create a configuration of your facility.

1. Run `./nkn-esi coordination-node init configs/facility`.
2. Run `/.nkn-esi coordination-node start configs/facility.json configs/facility.secret`.

You should now be at a shell terminal.

#### Exchange

You must now create a configuration of your exchange:

1. Run `./nkn-esi coordination-node init configs/exchange`.
2. Run `/.nkn-esi coordination-node start configs/exchange.json configs/exchange.secret`.

3. You should now be at a shell terminal.

### Signing Up

To sign your exchange up to a registry, you must run `registry signup`, to which you will be prompted for the registry
public key. If you ever forget the public key for either your facility or exchange, you can see it in `configs/x.json`,
or print it out in the terminal with `info public`.

After signing up to the registry, you will see that the registry has taken your key.

By signing up to a registry, you can now quickly and easily receive potential facilities.

### Querying

Your facility now needs to find the exchange. To do this, you can run `registry query` on the facility shell. It will
then prompt you for the registry public key and a country. This corresponds to the country listed in your exchange
config file. For now, just leave the default, either by pressing ENTER or typing in 'DC'.

You will then receive a new exchange. To see the list of exchanges you have found, type in `registry list`.

In the ESI, querying allows you to find exchanges based on shared or relevant details that may be important to you.

### Registering

Your facility now needs to register with the exchange. Run `facility request` in the facility shell, entering the public
key of the exchange. This will have the exchange send you a custom registration form. This registration form is totally
created and sent by the exchange, meaning that you can send whatever questions that may be important to your business.

To view the forms that you have, you can type `facility forms`.

To register, you must now run `facility register` in the facility shell, entering the public key of the exchange. This
will prompt you to answer a simple yes or no question. Either enter in Y, or just leave it and press ENTER to go with
the default.

### Creating Price Maps and Characteristics

The ESI describes the transaction process between a facility and an exchange.

1. The facility creates a price map and characteristics.
2. The exchange gets the facility's price map and characteristics.
3. The exchange proposes an offer to the facility with their own price map.
4. The facility responds to this offer, either with acceptance or their own counter offer.

