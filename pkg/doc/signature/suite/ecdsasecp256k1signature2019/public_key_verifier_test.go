/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ecdsasecp256k1signature2019

import (
	"testing"

	gojose "github.com/square/go-jose/v3"
	"github.com/stretchr/testify/require"

	"github.com/markcryptohash/aries-framework-go/pkg/crypto/tinkcrypto"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/jose/jwk"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/util/signature"
	kmsapi "github.com/markcryptohash/aries-framework-go/pkg/kms"
	"github.com/markcryptohash/aries-framework-go/pkg/kms/localkms"
	mockkms "github.com/markcryptohash/aries-framework-go/pkg/mock/kms"
	"github.com/markcryptohash/aries-framework-go/pkg/mock/storage"
	"github.com/markcryptohash/aries-framework-go/pkg/secretlock/noop"
)

func TestPublicKeyVerifier_Verify(t *testing.T) {
	signer, err := newCryptoSigner(kmsapi.ECDSASecp256k1TypeIEEEP1363)
	require.NoError(t, err)

	msg := []byte("test message")

	msgSig, err := signer.Sign(msg)
	require.NoError(t, err)

	pubKey := &verifier.PublicKey{
		Type: "EcdsaSecp256k1VerificationKey2019",

		JWK: &jwk.JWK{
			JSONWebKey: gojose.JSONWebKey{
				Algorithm: "ES256K",
				Key:       signer.PublicKey(),
			},
			Crv: "secp256k1",
			Kty: "EC",
		},
	}

	v := NewPublicKeyVerifier()

	err = v.Verify(pubKey, msg, msgSig)
	require.NoError(t, err)

	pubKey = &verifier.PublicKey{
		Type:  "EcdsaSecp256k1VerificationKey2019",
		Value: signer.PublicKeyBytes(),
	}

	err = v.Verify(pubKey, msg, msgSig)
	require.NoError(t, err)
}

func newCryptoSigner(keyType kmsapi.KeyType) (signature.Signer, error) {
	p, err := mockkms.NewProviderForKMS(storage.NewMockStoreProvider(), &noop.NoLock{})
	if err != nil {
		return nil, err
	}

	localKMS, err := localkms.New("local-lock://custom/master/key/", p)
	if err != nil {
		return nil, err
	}

	tinkCrypto, err := tinkcrypto.New()
	if err != nil {
		return nil, err
	}

	return signature.NewCryptoSigner(tinkCrypto, localKMS, keyType)
}
