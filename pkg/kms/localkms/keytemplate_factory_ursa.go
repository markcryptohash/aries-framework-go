//go:build ursa
// +build ursa

/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package localkms

import (
	tinkpb "github.com/google/tink/go/proto/tink_go_proto"

	clbld "github.com/markcryptohash/aries-framework-go/pkg/crypto/tinkcrypto/primitive/cl/blinder"
	clsgn "github.com/markcryptohash/aries-framework-go/pkg/crypto/tinkcrypto/primitive/cl/signer"
	"github.com/markcryptohash/aries-framework-go/pkg/kms"
)

// getKeyTemplate returns tink KeyTemplate associated with the provided keyType.
func getKeyTemplate(keyType kms.KeyType, opts ...kms.KeyOpts) (*tinkpb.KeyTemplate, error) {
	switch keyType {
	case kms.CLCredDefType:
		keyOpts := kms.NewKeyOpt()

		for _, opt := range opts {
			opt(keyOpts)
		}

		return clsgn.CredDefKeyTemplate(keyOpts.Attrs()), nil
	case kms.CLMasterSecretType:
		return clbld.MasterSecretKeyTemplate(), nil
	default:
		return keyTemplate(keyType, opts...)
	}
}
