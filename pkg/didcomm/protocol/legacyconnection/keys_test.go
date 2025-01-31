/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package legacyconnection

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/mediator"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/did"
	"github.com/markcryptohash/aries-framework-go/pkg/kms"
	"github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol"
	mockroute "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/mediator"
	mockstorage "github.com/markcryptohash/aries-framework-go/pkg/mock/storage"
)

func TestCreateNewKeyAndVM(t *testing.T) {
	k := newKMS(t, mockstorage.NewMockStoreProvider())

	p, err := New(&protocol.MockProvider{
		ServiceMap: map[string]interface{}{
			mediator.Coordination: &mockroute.MockMediatorSvc{},
		},
		CustomKMS: k,
	})
	require.NoError(t, err)

	t.Run("createNewKeyAndVM success", func(t *testing.T) {
		didDoc := &did.Doc{}

		p.ctx.keyType = kms.ED25519
		p.ctx.keyAgreementType = kms.X25519ECDHKWType

		err = p.ctx.createNewKeyAndVM(didDoc)
		require.NoError(t, err)
		require.Equal(t, ed25519VerificationKey2018, didDoc.VerificationMethod[0].Type)
		require.Equal(t, x25519KeyAgreementKey2019, didDoc.KeyAgreement[0].VerificationMethod.Type)
	})

	t.Run("createNewKeyAndVM invalid keyType export signing key", func(t *testing.T) {
		didDoc := &did.Doc{}

		p.ctx.keyType = kms.HMACSHA256Tag256Type // invalid signing key
		p.ctx.keyAgreementType = kms.X25519ECDHKWType

		err = p.ctx.createNewKeyAndVM(didDoc)
		require.EqualError(t, err, "createSigningVM: createAndExportPubKeyBytes: failed to export new public key bytes: "+
			"exportPubKeyBytes: failed to export marshalled key: exportPubKeyBytes: failed to get public keyset "+
			"handle: keyset.Handle: keyset.Handle: keyset contains a non-private key")
		require.Empty(t, didDoc.VerificationMethod)
		require.Empty(t, didDoc.KeyAgreement)
	})
}

func TestCreateSigningVM(t *testing.T) {
	k := newKMS(t, mockstorage.NewMockStoreProvider())

	p, err := New(&protocol.MockProvider{
		ServiceMap: map[string]interface{}{
			mediator.Coordination: &mockroute.MockMediatorSvc{},
		},
		CustomKMS: k,
	})
	require.NoError(t, err)

	t.Run("createSigningVM success", func(t *testing.T) {
		p.ctx.keyType = kms.ED25519

		svm, err := p.ctx.createSigningVM()
		require.NoError(t, err)
		require.NotEmpty(t, svm)
	})

	t.Run("createSigningVM with empty vmType", func(t *testing.T) {
		p.ctx.keyType = ""

		svm, err := p.ctx.createSigningVM()
		require.EqualError(t, err, "createSigningVM: createAndExportPubKeyBytes: failed to create new key: "+
			"failed to create new key, missing key type")
		require.Empty(t, svm)
	})

	t.Run("createSigningVM with unsupported keyType", func(t *testing.T) {
		p.ctx.keyType = kms.X25519ECDHKW

		svm, err := p.ctx.createSigningVM()
		require.EqualError(t, err, "createSigningVM: unsupported verification method: 'X25519KeyAgreementKey2019'")
		require.Empty(t, svm)
	})
}

func TestCreateEncryptionVM(t *testing.T) {
	k := newKMS(t, mockstorage.NewMockStoreProvider())

	p, err := New(&protocol.MockProvider{
		ServiceMap: map[string]interface{}{
			mediator.Coordination: &mockroute.MockMediatorSvc{},
		},
		CustomKMS: k,
	})
	require.NoError(t, err)

	t.Run("createEncryptionVM success", func(t *testing.T) {
		p.ctx.keyAgreementType = kms.X25519ECDHKW

		evm, err := p.ctx.createEncryptionVM()
		require.NoError(t, err)
		require.NotEmpty(t, evm)
	})

	t.Run("createEncryptionVM with empty keyAgreementType", func(t *testing.T) {
		p.ctx.keyAgreementType = ""

		evm, err := p.ctx.createEncryptionVM()
		require.EqualError(t, err, "createEncryptionVM: createAndExportPubKeyBytes: failed to create new key: "+
			"failed to create new key, missing key type")
		require.Empty(t, evm)
	})

	t.Run("createEncryptionVM with unsupported keyType", func(t *testing.T) {
		p.ctx.keyAgreementType = kms.ED25519Type

		evm, err := p.ctx.createEncryptionVM()
		require.EqualError(t, err, "unsupported verification method for KeyAgreement: 'Ed25519VerificationKey2018'")
		require.Empty(t, evm)
	})
}
