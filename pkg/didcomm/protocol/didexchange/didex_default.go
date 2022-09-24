// +build !ACAPyInterop

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package didexchange

import (
	"fmt"

	"github.com/markcryptohash/aries-framework-go/pkg/doc/did"
)

const (
	doACAPyInterop = false
)

// Interop: this is a stub method, that is substituted for special functionality
//   when the ACAPyInterop flag is enabled.
//   This can be removed when https://github.com/markcryptohash/aries-cloudagent-python/issues/1048 is fixed.
func convertPeerToSov(doc *did.Doc) (*did.Doc, error) {
	return doc, nil
}

func interopRecipientKey(doc *did.Doc) (string, error) {
	return "", fmt.Errorf("recipientKeyAsDIDKey: invalid DID Doc service type: '%v'", doc.Service[0].Type)
}
