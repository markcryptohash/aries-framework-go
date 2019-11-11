[![Release](https://img.shields.io/github/release/hyperledger/aries-framework-go.svg?style=flat-square)](https://github.com/hyperledger/aries-framework-go/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/trustbloc/aries-framework-go/master/LICENSE)
[![Godocs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hyperledger/aries-framework-go)

[![Build Status](https://github.com/hyperledger/aries-framework-go/workflows/build/badge.svg)](https://github.com/hyperledger/aries-framework-go/actions)
[![codecov](https://codecov.io/gh/hyperledger/aries-framework-go/branch/master/graph/badge.svg)](https://codecov.io/gh/hyperledger/aries-framework-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperledger/aries-framework-go)](https://goreportcard.com/report/github.com/hyperledger/aries-framework-go)

# <p><img src="https://raw.githubusercontent.com/hyperledger/aries-rfcs/1371a4807ead74c36ea7d5af909064ec491b78c1/collateral/Hyperledger_Aries_Logo_Color.png" height="50px" alt="Hyperledger Aries">Framework Go</p>

Hyperledger Aries Framework Go enables trusted communication and data exchange based on interoperable distributed ledger technologies (DLTs) and peer-to-peer (P2P) interactions.
We provide a flexible toolkit to enable the usage of decentralized identifiers (DIDs), DID-to-DID communications, verifiable credential exchange, transaction authorizations, and data communication protocols. From these building blocks, implementors can build agents, mediators and other DIDComm features in a manner that is agnostic to a particular DID network or governance framework.

We aim to provide Go implementations of:

- Decentralized identity standards including [W3C decentralized identifiers](https://w3c.github.io/did-core/) (DIDs), [W3C DID resolution](https://w3c-ccg.github.io/did-resolution/), and [W3C verifiable credentials](https://w3c.github.io/vc-data-model/).
- Decentralized data communication protocols anchored in DIDs: [DIDComm](https://github.com/hyperledger/aries-rfcs/blob/master/concepts/0005-didcomm).
- A pluggable dependency framework, where implementors can customize primitives via Service Provider Interfaces (SPIs). We have a "batteries included" model where default primitives are included -- such as a key management system (KMS), crypto, data storage, digital hub integration, etc.

We aim to enable usage of our protocol implementations in a wide variety of edge and cloud environments including servers, browsers, mobile, and devices.
API bindings are supplied to enable these environments including:

- Go
- REST
- C (future)
- WebAssembly (future)

We implement demonstrations and test cases, that require a ledger system, using [DIF Sidetree protocol](https://github.com/decentralized-identity/sidetree/blob/master/docs/protocol.md) as this protocol enables generic decentralized ledger systems to operate as a DID network.

## Build
### Prerequisites (General)
- Go 1.13

### Prerequisites (for running tests and demos)
- Go 1.13
- Docker
- Docker-Compose
- Make

### Targets
```
# run all the project build targets
make all

# run linter checks
make checks

# run unit tests
make unit-test
```

## Table of Contents

- [Running bdd tests](docs/test/bdd_instructions.md)
- [Running demo locally](docs/demo/openapi_demo_instructions.md)
- [Running agent locally](docs/agentd/agent_CLI.md)
- [Running agent using docker](docs/agentd/agent_docker_instructions.md)
- [Agent webhook support](docs/agentd/agent_webhook.md)
- [Generating Controller REST API specifications](docs/spec/openapi_spec_instructions.md)

### Crypto material generation for tests
For unit-tests, crypto material is generated under:

`pkg/didcomm/transport/http/testdata`

using the `openssl` tool. 

It is generated automatically when running unit tests. 

If you wish to regenerate it, you can delete this folder and:
1. run `make unit-test`
 or
2. cd into `pkg/didcomm/transport/http/` and run `go generate`

### Verifiable Credential Test Suite
To test compatibility of the verifiable credential packages with 
[W3C Verifiable Claims Working Group Test Suite](https://github.com/w3c/vc-test-suite), run `make vc-test-suite`.
The result of the test suite is generated as `vc-test-suite/suite/implementations/aries-framework-go-report.json`.

### Documentation

Agent documentation can be viewed at [GoDoc](https://godoc.org/github.com/hyperledger/aries-framework-go).

The packages intended for end developer usage are within the pkg/client folder along with the main agent package (pkg/framework/aries).

## Contributing

Found a bug? Ready to submit a PR? Want to submit a proposal for your grand
idea? Follow our [guidelines](.github/CONTRIBUTING.md) for more information
to get you started!

## License

Hyperledger Aries Framework Go is licensed under the [Apache License Version 2.0 (Apache-2.0)](LICENSE).

Hyperledger Aries Framework Go [documentation](docs) is licensed under the [Creative Commons Attribution 4.0 International License (CC-BY-4.0)](http://creativecommons.org/licenses/by/4.0/).
