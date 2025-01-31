/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package issuecredential

import "github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"

// Handler describes middleware interface.
type Handler interface {
	Handle(metadata Metadata) error
}

// Middleware function receives next handler and returns handler that needs to be executed.
type Middleware func(next Handler) Handler

// HandlerFunc is a helper type which implements the middleware Handler interface.
type HandlerFunc func(metadata Metadata) error

// Handle implements function to satisfy the Handler interface.
func (hf HandlerFunc) Handle(metadata Metadata) error {
	return hf(metadata)
}

// Metadata provides helpful information for the processing.
type Metadata interface {
	// Message contains the original inbound/outbound message
	Message() service.DIDCommMsg
	// OfferCredentialV2 is pointer to the message provided by the user through the Continue function.
	OfferCredentialV2() *OfferCredentialV2
	// ProposeCredentialV2 is pointer to the message provided by the user through the Continue function.
	ProposeCredentialV2() *ProposeCredentialV2
	// IssueCredential is pointer to the message provided by the user through the Continue function.
	IssueCredentialV2() *IssueCredentialV2
	// RequestCredential is pointer to message provided by the user through the Continue function.
	RequestCredentialV2() *RequestCredentialV2
	// CredentialNames is a slice which contains credential names provided by the user through the Continue function.
	CredentialNames() []string
	// StateName provides the state name
	StateName() string
	// Properties provides the possibility to set properties
	Properties() map[string]interface{}
}
