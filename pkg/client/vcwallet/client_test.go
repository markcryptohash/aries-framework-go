/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vcwallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/markcryptohash/aries-framework-go/internal/testdata"
	"github.com/markcryptohash/aries-framework-go/pkg/client/outofband"
	"github.com/markcryptohash/aries-framework-go/pkg/crypto/tinkcrypto"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/model"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/common/service"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/didexchange"
	issuecredentialsvc "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	"github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/mediator"
	outofbandSvc "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/outofband"
	oobv2 "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/outofbandv2"
	presentproofSvc "github.com/markcryptohash/aries-framework-go/pkg/didcomm/protocol/presentproof"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/did"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/presexch"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/util"
	"github.com/markcryptohash/aries-framework-go/pkg/doc/verifiable"
	vdrapi "github.com/markcryptohash/aries-framework-go/pkg/framework/aries/api/vdr"
	mockoutofbandv2 "github.com/markcryptohash/aries-framework-go/pkg/internal/gomocks/client/outofbandv2"
	"github.com/markcryptohash/aries-framework-go/pkg/internal/ldtestutil"
	"github.com/markcryptohash/aries-framework-go/pkg/kms"
	cryptomock "github.com/markcryptohash/aries-framework-go/pkg/mock/crypto"
	mockdidexchange "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/didexchange"
	mockissuecredential "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/issuecredential"
	mockmediator "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/mediator"
	mockoutofband "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/outofband"
	mockpresentproof "github.com/markcryptohash/aries-framework-go/pkg/mock/didcomm/protocol/presentproof"
	mockprovider "github.com/markcryptohash/aries-framework-go/pkg/mock/provider"
	"github.com/markcryptohash/aries-framework-go/pkg/mock/secretlock"
	mockstorage "github.com/markcryptohash/aries-framework-go/pkg/mock/storage"
	mockvdr "github.com/markcryptohash/aries-framework-go/pkg/mock/vdr"
	"github.com/markcryptohash/aries-framework-go/pkg/secretlock/local/masterlock/pbkdf2"
	"github.com/markcryptohash/aries-framework-go/pkg/store/connection"
	"github.com/markcryptohash/aries-framework-go/pkg/vdr/key"
	"github.com/markcryptohash/aries-framework-go/pkg/wallet"
	"github.com/markcryptohash/aries-framework-go/spi/storage"
)

const (
	samplePassPhrase        = "fakepassphrase"
	sampleRemoteKMSAuth     = "sample-auth-token"
	sampleKeyServerURL      = "sample/keyserver/test"
	sampleUserID            = "sample-user01"
	toBeImplementedErr      = "to be implemented"
	sampleClientErr         = "sample client err"
	sampleDIDKey            = "did:key:z6MknC1wwS6DEYwtGbZZo2QvjQjkh2qSBjb4GYmbye8dv4S5"
	sampleDIDKey2           = "did:key:z6MkwFKUCsf8wvn6eSSu1WFAKatN1yexiDM7bf7pZLSFjdz6"
	sampleInvalidDIDContent = `{"@context": ["https://w3id.org/did/v1"],"id": "did:example:sampleInvalidDIDContent"}`
	webRedirectStatusKey    = "status"
	webRedirectURLKey       = "url"
	exampleWebRedirect      = "http://example.com/sample"
	sampleMsgComment        = "sample mock msg"
)

func TestCreateProfile(t *testing.T) {
	t.Run("test create new wallet client using local kms passphrase", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)
		require.NoError(t, ProfileExists(sampleUserID, mockctx))

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test create new wallet client using local kms secret lock service", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := CreateProfile(sampleUserID, mockctx, wallet.WithSecretLockService(&secretlock.MockSecretLock{}))
		require.NoError(t, err)
		require.NoError(t, ProfileExists(sampleUserID, mockctx))

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test create new wallet client using remote kms key server URL", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.NoError(t, err)
		require.NoError(t, ProfileExists(sampleUserID, mockctx))

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test create new wallet failure", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := CreateProfile(sampleUserID, mockctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid create profile options")
		require.True(t, errors.Is(ProfileExists(sampleUserID, mockctx), wallet.ErrProfileNotFound))

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)
	})

	t.Run("test create new wallet failure - create store error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		mockctx.StorageProviderValue = &mockstorage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf(sampleClientErr),
		}

		err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)

		err = ProfileExists(sampleUserID, mockctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)
	})

	t.Run("test create new wallet failure - save profile error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		mockctx.StorageProviderValue = &mockstorage.MockStoreProvider{
			Store: &mockstorage.MockStore{
				ErrPut: fmt.Errorf(sampleClientErr),
			},
		}

		err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)
	})

	t.Run("test create new wallet failure - create content store error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		mockctx.StorageProviderValue = &mockStorageProvider{
			MockStoreProvider: mockstorage.NewMockStoreProvider(),
			failure:           fmt.Errorf(sampleClientErr),
		}

		err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)
		require.Empty(t, vcWallet)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("test update wallet client using local kms passphrase", func(t *testing.T) {
		mockctx := newMockProvider(t)
		createSampleProfile(t, mockctx)

		err := UpdateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test update wallet client using local kms secret lock service", func(t *testing.T) {
		mockctx := newMockProvider(t)
		createSampleProfile(t, mockctx)

		err := UpdateProfile(sampleUserID, mockctx, wallet.WithSecretLockService(&secretlock.MockSecretLock{}))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test update wallet client using remote kms key server URL", func(t *testing.T) {
		mockctx := newMockProvider(t)
		createSampleProfile(t, mockctx)

		err := UpdateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test update wallet failure", func(t *testing.T) {
		mockctx := newMockProvider(t)
		createSampleProfile(t, mockctx)

		err := UpdateProfile(sampleUserID, mockctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid create profile options")

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test update wallet failure - profile doesn't exists", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := UpdateProfile(sampleUserID, mockctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "profile does not exist")

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)
	})

	t.Run("test update wallet failure - create store error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		mockctx.StorageProviderValue = &mockstorage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf(sampleClientErr),
		}

		err := UpdateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)
	})

	t.Run("test update wallet failure - save profile error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		createSampleProfile(t, mockctx)

		mockctx.StorageProviderValue.(*mockstorage.MockStoreProvider).Store.ErrPut = fmt.Errorf(sampleClientErr)

		err := UpdateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})
}

func TestNew(t *testing.T) {
	t.Run("test get client", func(t *testing.T) {
		mockctx := newMockProvider(t)
		// create a wallet
		err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)
	})

	t.Run("test get client unlocked", func(t *testing.T) {
		mockctx := newMockProvider(t)
		// create a wallet
		err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		token, err := vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("test get client unlock failure - wrong passphrase", func(t *testing.T) {
		mockctx := newMockProvider(t)
		// create a wallet
		err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase+"ss"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "message authentication failed")
		require.Empty(t, vcWallet)
	})

	t.Run("test get client by invalid userID", func(t *testing.T) {
		mockctx := newMockProvider(t)
		err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUserID+"invalid", mockctx)
		require.Empty(t, vcWallet)
		require.Error(t, err)
		require.Contains(t, err.Error(), "profile does not exist")
	})

	t.Run("test update wallet failure - save profile error", func(t *testing.T) {
		mockctx := newMockProvider(t)
		mockctx.StorageProviderValue = &mockstorage.MockStoreProvider{
			ErrOpenStoreHandle: fmt.Errorf(sampleClientErr),
		}

		vcWallet, err := New(sampleUserID, mockctx)
		require.Error(t, err)
		require.Empty(t, vcWallet)
		require.Contains(t, err.Error(), sampleClientErr)
	})
}

func TestClient_OpenClose(t *testing.T) {
	t.Run("test open & close wallet using local kms passphrase", func(t *testing.T) {
		sampleUser := uuid.New().String()
		mockctx := newMockProvider(t)

		err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
		require.NoError(t, err)

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		// get token
		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		token, err := vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		defer vcWallet.Close()

		// try again
		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.Error(t, err)
		require.True(t, errors.Is(err, wallet.ErrAlreadyUnlocked))
		token, err = vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// close wallet
		require.True(t, vcWallet.Close())
		require.False(t, vcWallet.Close())

		// try to open with wrong passphrase
		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase + "wrong"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "message authentication failed")
		token, err = vcWallet.auth()
		require.Empty(t, token)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWalletLocked))
	})

	t.Run("test open & close wallet using secret lock service", func(t *testing.T) {
		sampleUser := uuid.New().String()
		mockctx := newMockProvider(t)

		masterLock, err := pbkdf2.NewMasterLock(samplePassPhrase, sha256.New, 0, nil)
		require.NoError(t, err)

		err = CreateProfile(sampleUser, mockctx, wallet.WithSecretLockService(masterLock))
		require.NoError(t, err)

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		// get token
		err = vcWallet.Open(wallet.WithUnlockBySecretLockService(masterLock))
		require.NoError(t, err)
		token, err := vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		defer vcWallet.Close()

		// try again
		err = vcWallet.Open(wallet.WithUnlockBySecretLockService(masterLock))
		require.Error(t, err)
		require.True(t, errors.Is(err, wallet.ErrAlreadyUnlocked))
		token, err = vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// close wallet
		require.True(t, vcWallet.Close())
		require.False(t, vcWallet.Close())

		// try to open with wrong secret lock service
		badLock, err := pbkdf2.NewMasterLock(samplePassPhrase+"wrong", sha256.New, 0, nil)
		require.NoError(t, err)

		err = vcWallet.Open(wallet.WithUnlockBySecretLockService(badLock))
		require.Error(t, err)
		require.Contains(t, err.Error(), "message authentication failed")
		token, err = vcWallet.auth()
		require.Empty(t, token)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWalletLocked))
	})

	t.Run("test open & close wallet using remote kms URL", func(t *testing.T) {
		sampleUser := uuid.New().String()
		mockctx := newMockProvider(t)

		err := CreateProfile(sampleUser, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
		require.NoError(t, err)

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		// get token
		err = vcWallet.Open(wallet.WithUnlockByAuthorizationToken(sampleRemoteKMSAuth))
		require.NoError(t, err)
		token, err := vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		defer vcWallet.Close()

		// try again
		err = vcWallet.Open(wallet.WithUnlockByAuthorizationToken(sampleRemoteKMSAuth))
		require.Error(t, err)
		require.True(t, errors.Is(err, wallet.ErrAlreadyUnlocked))
		token, err = vcWallet.auth()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// close wallet
		require.True(t, vcWallet.Close())
		require.False(t, vcWallet.Close())
		token, err = vcWallet.auth()
		require.Empty(t, token)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWalletLocked))
	})
}

func TestClient_Export(t *testing.T) {
	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx)
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	result, err := vcWalletClient.Export("")
	require.Empty(t, result)
	require.Error(t, err)
	require.EqualError(t, err, toBeImplementedErr)
}

func TestClient_Import(t *testing.T) {
	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx)
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	err = vcWalletClient.Import("", nil)
	require.Error(t, err)
	require.EqualError(t, err, toBeImplementedErr)
}

func TestClient_Add(t *testing.T) {
	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	err = vcWalletClient.Add(wallet.Metadata, testdata.SampleWalletContentMetadata)
	require.NoError(t, err)

	// try locked wallet
	vcWalletClient, err = New(sampleUserID, mockctx)
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	err = vcWalletClient.Add(wallet.Metadata, testdata.SampleWalletContentMetadata)
	require.Contains(t, err.Error(), "wallet locked")
}

func TestClient_Get(t *testing.T) {
	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	err = vcWalletClient.Add(wallet.Metadata, testdata.SampleWalletContentMetadata)
	require.NoError(t, err)

	content, err := vcWalletClient.Get(wallet.Metadata, "did:example:123456789abcdefghi")
	require.NoError(t, err)
	require.NotEmpty(t, content)
	require.Equal(t, string(testdata.SampleWalletContentMetadata), string(content))

	// try locked wallet
	require.True(t, vcWalletClient.Close())
	content, err = vcWalletClient.Get(wallet.Metadata, "did:example:123456789abcdefghi")
	require.True(t, errors.Is(err, ErrWalletLocked))
	require.Empty(t, content)
}

func TestClient_GetAll(t *testing.T) {
	const vcContent = `{
      "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.w3.org/2018/credentials/examples/v1"
      ],
      "id": "%s",
      "issuer": {
        "id": "did:example:76e12ec712ebc6f1c221ebfeb1f"
      },
      "type": [
        "VerifiableCredential",
        "UniversityDegreeCredential"
      ]
    }`

	const orgCollection = `{
                    "@context": ["https://w3id.org/wallet/v1"],
                    "id": "did:example:acme123456789abcdefghi",
                    "type": "Organization",
                    "name": "Acme Corp.",
                    "image": "https://via.placeholder.com/150",
                    "description" : "A software company.",
                    "tags": ["professional", "organization"],
                    "correlation": ["4058a72a-9523-11ea-bb37-0242ac130002"]
                }`

	const collectionID = "did:example:acme123456789abcdefghi"

	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	// save test data
	const count = 5

	for i := 0; i < count; i++ {
		require.NoError(t, vcWalletClient.Add(wallet.Credential, []byte(fmt.Sprintf(vcContent, uuid.New().String()))))
	}

	// save a collection
	require.NoError(t, vcWalletClient.Add(wallet.Collection, []byte(orgCollection)))

	// save contents by collection
	for i := 0; i < count; i++ {
		require.NoError(t, vcWalletClient.Add(wallet.Credential,
			[]byte(fmt.Sprintf(vcContent, uuid.New().String())), wallet.AddByCollection(collectionID)))
	}

	// get all by content
	vcs, err := vcWalletClient.GetAll(wallet.Credential)
	require.NoError(t, err)
	require.Len(t, vcs, count*2)

	// get all by content & collection
	vcs, err = vcWalletClient.GetAll(wallet.Credential, wallet.FilterByCollection(collectionID))
	require.NoError(t, err)
	require.Len(t, vcs, count)

	// try locked wallet
	require.True(t, vcWalletClient.Close())
	vcs, err = vcWalletClient.GetAll(wallet.Credential, wallet.FilterByCollection(collectionID))
	require.True(t, errors.Is(err, ErrWalletLocked))
	require.Empty(t, vcs)
}

func TestClient_Remove(t *testing.T) {
	mockctx := newMockProvider(t)
	err := CreateProfile(sampleUserID, mockctx, wallet.WithKeyServerURL(sampleKeyServerURL))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	err = vcWalletClient.Add(wallet.Metadata, testdata.SampleWalletContentMetadata)
	require.NoError(t, err)

	content, err := vcWalletClient.Get(wallet.Metadata, "did:example:123456789abcdefghi")
	require.NoError(t, err)
	require.NotEmpty(t, content)

	err = vcWalletClient.Remove(wallet.Metadata, "did:example:123456789abcdefghi")
	require.NoError(t, err)

	content, err = vcWalletClient.Get(wallet.Metadata, "did:example:123456789abcdefghi")
	require.Empty(t, content)
	require.Error(t, err)
	require.True(t, errors.Is(err, storage.ErrDataNotFound))

	// try locked wallet
	require.True(t, vcWalletClient.Close())
	err = vcWalletClient.Remove(wallet.Metadata, "did:example:123456789abcdefghi")
	require.True(t, errors.Is(err, ErrWalletLocked))
}

func TestClient_Query(t *testing.T) {
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()

				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}

				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		},
	}

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = customVDR
	mockctx.CryptoValue = &cryptomock.Crypto{}

	vc1, err := (&verifiable.Credential{
		Context: []string{verifiable.ContextURI},
		Types:   []string{verifiable.VCType},
		ID:      "http://example.edu/credentials/9999",
		CustomFields: map[string]interface{}{
			"first_name": "Jesse",
		},
		Issued: &util.TimeWrapper{
			Time: time.Now(),
		},
		Issuer: verifiable.Issuer{
			ID: "did:example:76e12ec712ebc6f1c221ebfeb1f",
		},
		Subject: uuid.New().String(),
	}).MarshalJSON()
	require.NoError(t, err)

	sampleVC := strings.ReplaceAll(string(testdata.SampleUDCVCWithCredentialSchema),
		"https://example.com/schema", verifiable.ContextURI)
	vcForQuery := []byte(strings.ReplaceAll(sampleVC,
		"http://example.edu/credentials/1872", "http://example.edu/credentials/1879"))
	vcForDerive := testdata.SampleUDCVCWithProofBBS

	err = CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NotEmpty(t, vcWalletClient)
	require.NoError(t, err)

	require.NoError(t, vcWalletClient.Add(wallet.Credential, vc1))
	require.NoError(t, vcWalletClient.Add(wallet.Credential, vcForQuery))
	require.NoError(t, vcWalletClient.Add(wallet.Credential, vcForDerive))

	pd := &presexch.PresentationDefinition{
		ID: uuid.New().String(),
		InputDescriptors: []*presexch.InputDescriptor{{
			ID: uuid.New().String(),
			Schema: []*presexch.Schema{{
				URI: fmt.Sprintf("%s#%s", verifiable.ContextID, verifiable.VCType),
			}},
			Constraints: &presexch.Constraints{
				Fields: []*presexch.Field{{
					Path: []string{"$.first_name"},
				}},
			},
		}},
	}

	// presentation exchange
	pdJSON, err := json.Marshal(pd)
	require.NoError(t, err)
	require.NotEmpty(t, pdJSON)

	// query by example
	queryByExample := []byte(strings.ReplaceAll(string(testdata.SampleWalletQueryByExample),
		"did:example:abcd", verifiable.ContextURI))
	// query by frame
	queryByFrame := testdata.SampleWalletQueryByFrame

	t.Run("test wallet queries", func(t *testing.T) {
		tests := []struct {
			name        string
			params      []*wallet.QueryParams
			resultCount int
			vcCount     map[int]int
			error       string
		}{
			{
				name: "query by presentation exchange - success",
				params: []*wallet.QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
				},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by example - success",
				params: []*wallet.QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
				},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by frame - success",
				params: []*wallet.QueryParams{
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
				},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "did authorization - success",
				params: []*wallet.QueryParams{
					{Type: "DIDAuth"},
				},
				resultCount: 1,
				vcCount:     map[int]int{0: 0},
			},
			{
				name: "multiple queries - success",
				params: []*wallet.QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
				},
				resultCount: 2,
				vcCount:     map[int]int{0: 1, 1: 2},
			},
			{
				name: "invalid query type",
				params: []*wallet.QueryParams{
					{Type: "invalid"},
				},
				error: "unsupported query type",
			},
			{
				name:   "empty query type",
				params: []*wallet.QueryParams{},
				error:  "no result found",
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.name, func(t *testing.T) {
				results, err := vcWalletClient.Query(tc.params...)

				if tc.error != "" {
					require.Empty(t, results)
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.error)

					return
				}

				require.NoError(t, err)
				require.Len(t, results, tc.resultCount)

				for i, result := range results {
					require.Len(t, result.Credentials(), tc.vcCount[i])
				}
			})
		}

		// test wallet locked
		require.True(t, vcWalletClient.Close())

		results, err := vcWalletClient.Query(&wallet.QueryParams{
			Type: "DIDAuth",
		})
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, results)
	})
}

func TestClient_Issue(t *testing.T) {
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()

				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}

				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		},
	}

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = customVDR
	mockctx.CryptoValue = &cryptomock.Crypto{}

	err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("Test VC wallet client issue using controller - success", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// save a DID & corresponding key
		require.NoError(t, vcWalletClient.Add(wallet.Key, testdata.SampleWalletContentKeyBase58))
		require.NoError(t, vcWalletClient.Add(wallet.DIDResolutionResponse, testdata.SampleDocResolutionResponse))

		result, err := vcWalletClient.Issue(testdata.SampleUDCVC, &wallet.ProofOptions{
			Controller: sampleDIDKey,
		})

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Proofs)
	})

	t.Run("Test VC wallet client issue using controller - failure", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// sign with just controller
		result, err := vcWalletClient.Issue(testdata.SampleUDCVC, &wallet.ProofOptions{
			Controller: sampleDIDKey2,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read json keyset from reader")
		require.Empty(t, result)
	})

	t.Run("Test VC wallet client issue using controller - wallet locked", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx)
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// sign with just controller
		result, err := vcWalletClient.Issue(testdata.SampleUDCVC, &wallet.ProofOptions{
			Controller: sampleDIDKey,
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, result)
	})
}

func TestClient_Prove(t *testing.T) {
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()

				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}

				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		},
	}

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = customVDR
	mockctx.CryptoValue = &cryptomock.Crypto{}

	err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("Test VC wallet client prove using controller - success", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// save a credential, DID & key
		require.NoError(t, vcWalletClient.Add(wallet.Credential, testdata.SampleUDCVC))
		require.NoError(t, vcWalletClient.Add(wallet.Key, testdata.SampleWalletContentKeyBase58))
		require.NoError(t, vcWalletClient.Add(wallet.DIDResolutionResponse, testdata.SampleDocResolutionResponse))

		result, err := vcWalletClient.Prove(&wallet.ProofOptions{Controller: sampleDIDKey},
			wallet.WithStoredCredentialsToProve("http://example.edu/credentials/1872"),
			wallet.WithRawCredentialsToProve(testdata.SampleUDCVC),
		)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Proofs)
	})

	t.Run("Test VC wallet client prove using controller - failure", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		require.NoError(t, vcWalletClient.Remove(wallet.Credential, "http://example.edu/credentials/1872"))
		require.NoError(t, vcWalletClient.Add(wallet.Credential, testdata.SampleUDCVC))

		result, err := vcWalletClient.Prove(&wallet.ProofOptions{Controller: sampleDIDKey2},
			wallet.WithStoredCredentialsToProve("http://example.edu/credentials/1872"),
			wallet.WithRawCredentialsToProve(testdata.SampleUDCVC),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read json keyset from reader")
		require.Empty(t, result)
	})

	t.Run("Test VC wallet client prove using controller - wallet locked", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		require.NoError(t, vcWalletClient.Remove(wallet.Credential, "http://example.edu/credentials/1872"))
		require.NoError(t, vcWalletClient.Add(wallet.Credential, testdata.SampleUDCVC))

		vcWalletClient.Close()

		result, err := vcWalletClient.Prove(&wallet.ProofOptions{Controller: sampleDIDKey},
			wallet.WithStoredCredentialsToProve("http://example.edu/credentials/1872"),
			wallet.WithRawCredentialsToProve(testdata.SampleUDCVC),
		)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, result)
	})
}

func TestClient_Verify(t *testing.T) {
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()

				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}

				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		},
	}

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = customVDR
	mockctx.CryptoValue = &cryptomock.Crypto{}

	err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("Test VC wallet verify credential - success", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// store credential in wallet
		require.NoError(t, vcWalletClient.Add(wallet.Credential, testdata.SampleUDCVCWithProof))

		// verify stored VC
		ok, err := vcWalletClient.Verify(wallet.WithStoredCredentialToVerify("http://example.edu/credentials/1872"))
		require.NoError(t, err)
		require.True(t, ok)

		// verify raw VC
		ok, err = vcWalletClient.Verify(wallet.WithRawCredentialToVerify(testdata.SampleUDCVCWithProof))
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("Test VC wallet verify credential - invalid signature", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// store tampered credential in wallet
		tamperedVC := strings.ReplaceAll(string(testdata.SampleUDCVCWithProof),
			`"name": "Example University"`, `"name": "Fake University"`)
		require.NoError(t, vcWalletClient.Remove(wallet.Credential, "http://example.edu/credentials/1872"))
		require.NoError(t, vcWalletClient.Add(wallet.Credential, []byte(tamperedVC)))

		ok, err := vcWalletClient.Verify(wallet.WithStoredCredentialToVerify("http://example.edu/credentials/1872"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
		require.False(t, ok)

		ok, err = vcWalletClient.Verify(wallet.WithRawCredentialToVerify([]byte(tamperedVC)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
		require.False(t, ok)
	})

	t.Run("Test VC wallet verify presentation - success", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		// verify raw VC
		ok, err := vcWalletClient.Verify(wallet.WithRawPresentationToVerify(testdata.SampleUDCPresentation))
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("Test VC wallet verify presentation - invalid signature", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		defer vcWalletClient.Close()

		tamperedVP := strings.ReplaceAll(string(testdata.SampleUDCPresentation),
			`"holder": "did:key:z6MknC1wwS6DEYwtGbZZo2QvjQjkh2qSBjb4GYmbye8dv4S5"`,
			`"holder": "did:key:z6MknC1wwS6DEYwtGbZZo2QvjQjkh2qSBjb4GYmbye8dv464"`)

		ok, err := vcWalletClient.Verify(wallet.WithRawPresentationToVerify([]byte(tamperedVP)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature")
		require.False(t, ok)
	})

	t.Run("Test VC wallet verify presentation - test wallet locked", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx)
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		ok, err := vcWalletClient.Verify(wallet.WithRawPresentationToVerify(testdata.SampleUDCPresentation))
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.False(t, ok)
	})
}

func TestWallet_Derive(t *testing.T) {
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()

				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}

				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		},
	}

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = customVDR

	customCrypto, err := tinkcrypto.New()
	require.NoError(t, err)

	mockctx.CryptoValue = customCrypto

	// create profile
	err = CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	// prepare frame
	var frameDoc map[string]interface{}

	require.NoError(t, json.Unmarshal(testdata.SampleFrame, &frameDoc))

	t.Run("Test derive a credential from wallet - success", func(t *testing.T) {
		walletInstance, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		require.NotEmpty(t, walletInstance)

		// save BBS VC in store
		require.NoError(t, walletInstance.Add(wallet.Credential, testdata.SampleUDCVCWithProofBBS))

		sampleNonce := uuid.New().String()

		verifyBBSProof := func(proofs []verifiable.Proof) {
			require.Len(t, proofs, 1)
			require.NotEmpty(t, proofs[0])
			require.Equal(t, proofs[0]["type"], "BbsBlsSignatureProof2020")
			require.NotEmpty(t, proofs[0]["nonce"])
			require.EqualValues(t, proofs[0]["nonce"], base64.StdEncoding.EncodeToString([]byte(sampleNonce)))
			require.NotEmpty(t, proofs[0]["proofValue"])
		}

		// derive stored credential
		vc, err := walletInstance.Derive(wallet.FromStoredCredential("http://example.edu/credentials/1872"),
			&wallet.DeriveOptions{
				Nonce: sampleNonce,
				Frame: frameDoc,
			})
		require.NoError(t, err)
		require.NotEmpty(t, vc)
		verifyBBSProof(vc.Proofs)

		// derive raw credential
		vc, err = walletInstance.Derive(wallet.FromRawCredential(testdata.SampleUDCVCWithProofBBS), &wallet.DeriveOptions{
			Nonce: sampleNonce,
			Frame: frameDoc,
		})
		require.NoError(t, err)
		require.NotEmpty(t, vc)
		verifyBBSProof(vc.Proofs)

		// derive from credential instance
		pkFetcher := verifiable.WithPublicKeyFetcher(
			verifiable.NewVDRKeyResolver(customVDR).PublicKeyFetcher(),
		)

		loader, err := ldtestutil.DocumentLoader()
		require.NoError(t, err)

		credential, err := verifiable.ParseCredential(testdata.SampleUDCVCWithProofBBS, pkFetcher,
			verifiable.WithJSONLDDocumentLoader(loader))
		require.NoError(t, err)
		vc, err = walletInstance.Derive(wallet.FromCredential(credential), &wallet.DeriveOptions{
			Nonce: sampleNonce,
			Frame: frameDoc,
		})
		require.NoError(t, err)
		require.NotEmpty(t, vc)
		verifyBBSProof(vc.Proofs)
	})

	t.Run("Test derive credential failures", func(t *testing.T) {
		walletInstance, err := New(sampleUserID, mockctx, wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NotEmpty(t, walletInstance)
		require.NoError(t, err)

		// invalid request
		vc, err := walletInstance.Derive(wallet.FromStoredCredential(""), &wallet.DeriveOptions{})
		require.Empty(t, vc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid request to derive credential")

		// credential not found in store
		vc, err = walletInstance.Derive(wallet.FromStoredCredential("invalid-id"), &wallet.DeriveOptions{})
		require.Empty(t, vc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "data not found")

		// invalid credential in store
		require.NoError(t, walletInstance.Add(wallet.Credential, []byte(sampleInvalidDIDContent)))

		vc, err = walletInstance.Derive(wallet.FromStoredCredential("did:example:sampleInvalidDIDContent"),
			&wallet.DeriveOptions{})
		require.Empty(t, vc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "credential type of unknown structure")

		// invalid raw credential
		vc, err = walletInstance.Derive(wallet.FromRawCredential([]byte(sampleInvalidDIDContent)), &wallet.DeriveOptions{})
		require.Empty(t, vc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "credential type of unknown structure")

		// try deriving wrong proof type - no BbsBlsSignature2020 proof present
		vc, err = walletInstance.Derive(wallet.FromRawCredential(testdata.SampleUDCVCWithProof), &wallet.DeriveOptions{
			Frame: frameDoc,
		})
		require.Empty(t, vc)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no BbsBlsSignature2020 proof present")
	})

	t.Run("Test derive credential failures - test wallet locked", func(t *testing.T) {
		vcWalletClient, err := New(sampleUserID, mockctx)
		require.NotEmpty(t, vcWalletClient)
		require.NoError(t, err)

		result, err := vcWalletClient.Derive(wallet.FromRawCredential(testdata.SampleUDCVCWithProof), &wallet.DeriveOptions{
			Frame: frameDoc,
		})
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, result)
	})
}

func TestClient_CreateKeyPair(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	vcWallet, err := New(sampleUser, mockctx)
	require.NoError(t, err)
	require.NotEmpty(t, vcWallet)

	err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("test creating key pair", func(t *testing.T) {
		keyPair, err := vcWallet.CreateKeyPair(kms.ED25519)
		require.NoError(t, err)
		require.NotEmpty(t, keyPair)
		require.NotEmpty(t, keyPair.KeyID)
		require.NotEmpty(t, keyPair.PublicKey)
	})

	t.Run("test failure while creating key pair", func(t *testing.T) {
		keyPair, err := vcWallet.CreateKeyPair(kms.KeyType("invalid"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create new key")
		require.Empty(t, keyPair)
	})

	t.Run("test failure while creating key pair (closed wallet)", func(t *testing.T) {
		require.True(t, vcWallet.Close())

		keyPair, err := vcWallet.CreateKeyPair(kms.KeyType("invalid"))
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, keyPair)
	})
}

func TestClient_Connect(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("test did connect success", func(t *testing.T) {
		sampleConnID := uuid.New().String()

		oobSvc := &mockoutofband.MockOobService{
			AcceptInvitationHandle: func(*outofbandSvc.Invitation, outofbandSvc.Options) (string, error) {
				return sampleConnID, nil
			},
		}
		mockctx.ServiceMap[outofbandSvc.Name] = oobSvc

		didexSvc := &mockdidexchange.MockDIDExchangeSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				ch <- service.StateMsg{
					Type:       service.PostState,
					StateID:    didexchange.StateIDCompleted,
					Properties: &mockdidexchange.MockEventProperties{ConnID: sampleConnID},
				}

				return nil
			},
		}
		mockctx.ServiceMap[didexchange.DIDExchange] = didexSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		connectionID, err := vcWallet.Connect(&outofband.Invitation{})
		require.NoError(t, err)
		require.Equal(t, sampleConnID, connectionID)
	})

	t.Run("test did connect failure", func(t *testing.T) {
		didexSvc := &mockdidexchange.MockDIDExchangeSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				return fmt.Errorf(sampleClientErr)
			},
		}
		mockctx.ServiceMap[didexchange.DIDExchange] = didexSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		connectionID, err := vcWallet.Connect(&outofband.Invitation{})
		require.Error(t, err)
		require.Contains(t, err.Error(), sampleClientErr)
		require.Empty(t, connectionID)
	})

	t.Run("test did connect failure - auth error", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		connectionID, err := vcWallet.Connect(&outofband.Invitation{})
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, connectionID)
	})
}

func TestClient_ProposePresentation(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	const (
		myDID    = "did:mydid:123"
		theirDID = "did:theirdid:123"
	)

	t.Run("test propose presentation success", func(t *testing.T) {
		sampleConnID := uuid.New().String()

		oobSvc := &mockoutofband.MockOobService{
			AcceptInvitationHandle: func(*outofbandSvc.Invitation, outofbandSvc.Options) (string, error) {
				return sampleConnID, nil
			},
		}
		mockctx.ServiceMap[outofbandSvc.Name] = oobSvc

		didexSvc := &mockdidexchange.MockDIDExchangeSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				ch <- service.StateMsg{
					Type:       service.PostState,
					StateID:    didexchange.StateIDCompleted,
					Properties: &mockdidexchange.MockEventProperties{ConnID: sampleConnID},
				}

				return nil
			},
		}
		mockctx.ServiceMap[didexchange.DIDExchange] = didexSvc

		thID := uuid.New().String()

		ppSvc := &mockpresentproof.MockPresentProofSvc{
			ActionsFunc: func() ([]presentproofSvc.Action, error) {
				return []presentproofSvc.Action{
					{
						PIID: thID,
						Msg: service.NewDIDCommMsgMap(&presentproofSvc.RequestPresentationV2{
							Comment: "mock msg",
						}),
						MyDID:    myDID,
						TheirDID: theirDID,
					},
				}, nil
			},
			HandleFunc: func(service.DIDCommMsg) (string, error) {
				return thID, nil
			},
		}
		mockctx.ServiceMap[presentproofSvc.Name] = ppSvc

		store, err := mockctx.StorageProvider().OpenStore(connection.Namespace)
		require.NoError(t, err)

		record := &connection.Record{
			ConnectionID: sampleConnID,
			MyDID:        myDID,
			TheirDID:     theirDID,
		}
		recordBytes, err := json.Marshal(record)
		require.NoError(t, err)
		require.NoError(t, store.Put(fmt.Sprintf("conn_%s", sampleConnID), recordBytes))

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		msg, err := vcWallet.ProposePresentation(&wallet.GenericInvitation{})
		require.NoError(t, err)
		require.NotEmpty(t, msg)
	})

	t.Run("test propose presentation failure", func(t *testing.T) {
		oobSvc := &mockoutofband.MockOobService{
			AcceptInvitationHandle: func(*outofbandSvc.Invitation, outofbandSvc.Options) (string, error) {
				return "", fmt.Errorf(sampleClientErr)
			},
		}
		mockctx.ServiceMap[outofbandSvc.Name] = oobSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		msg, err := vcWallet.ProposePresentation(&wallet.GenericInvitation{})
		require.Error(t, err)
		require.Empty(t, msg)
	})

	t.Run("test propose presentation failure - auth error", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		msg, err := vcWallet.ProposePresentation(&wallet.GenericInvitation{})
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, msg)
	})
}

func TestClient_PresentProof(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("test present proof success", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		response, err := vcWallet.PresentProof(uuid.New().String(), wallet.FromPresentation(&verifiable.Presentation{}))
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, model.AckStatusPENDING, response.Status)
	})

	t.Run("test present proof success - wait for done", func(t *testing.T) {
		thID := uuid.New().String()
		mockPresentProofSvc := &mockpresentproof.MockPresentProofSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				ch <- service.StateMsg{
					Type:    service.PostState,
					StateID: presentproofSvc.StateNameDone,
					Properties: &mockdidexchange.MockEventProperties{
						Properties: map[string]interface{}{
							webRedirectStatusKey: model.AckStatusOK,
							webRedirectURLKey:    exampleWebRedirect,
						},
					},
					Msg: &mockMsg{thID: thID},
				}

				return nil
			},
		}
		mockctx.ServiceMap[presentproofSvc.Name] = mockPresentProofSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		response, err := vcWallet.PresentProof(thID, wallet.FromPresentation(&verifiable.Presentation{}),
			wallet.WaitForDone(1*time.Millisecond))
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, model.AckStatusOK, response.Status)
		require.Equal(t, exampleWebRedirect, response.RedirectURL)
	})

	t.Run("test present proof failure - auth error", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		response, err := vcWallet.PresentProof(uuid.New().String(), wallet.FromPresentation(&verifiable.Presentation{}))
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, response)
	})
}

func TestClient_ProposeCredential(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	const (
		myDID    = "did:mydid:123"
		theirDID = "did:theirdid:123"
	)

	t.Run("test propose credential success", func(t *testing.T) {
		sampleConnID := uuid.New().String()

		oobSvc := &mockoutofband.MockOobService{
			AcceptInvitationHandle: func(*outofbandSvc.Invitation, outofbandSvc.Options) (string, error) {
				return sampleConnID, nil
			},
		}
		mockctx.ServiceMap[outofbandSvc.Name] = oobSvc

		didexSvc := &mockdidexchange.MockDIDExchangeSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				ch <- service.StateMsg{
					Type:       service.PostState,
					StateID:    didexchange.StateIDCompleted,
					Properties: &mockdidexchange.MockEventProperties{ConnID: sampleConnID},
				}

				return nil
			},
		}
		mockctx.ServiceMap[didexchange.DIDExchange] = didexSvc

		thID := uuid.New().String()

		icSvc := &mockissuecredential.MockIssueCredentialSvc{
			ActionsFunc: func() ([]issuecredentialsvc.Action, error) {
				return []issuecredentialsvc.Action{
					{
						PIID: thID,
						Msg: service.NewDIDCommMsgMap(&issuecredentialsvc.OfferCredentialV2{
							Comment: sampleMsgComment,
						}),
						MyDID:    myDID,
						TheirDID: theirDID,
					},
				}, nil
			},
			HandleFunc: func(service.DIDCommMsg) (string, error) {
				return thID, nil
			},
		}
		mockctx.ServiceMap[issuecredentialsvc.Name] = icSvc

		store, err := mockctx.StorageProvider().OpenStore(connection.Namespace)
		require.NoError(t, err)

		record := &connection.Record{
			ConnectionID: sampleConnID,
			MyDID:        myDID,
			TheirDID:     theirDID,
		}
		recordBytes, err := json.Marshal(record)
		require.NoError(t, err)
		require.NoError(t, store.Put(fmt.Sprintf("conn_%s", sampleConnID), recordBytes))

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		msg, err := vcWallet.ProposeCredential(&wallet.GenericInvitation{})
		require.NoError(t, err)
		require.NotEmpty(t, msg)

		offer := &issuecredentialsvc.OfferCredentialV2{}

		err = msg.Decode(offer)
		require.NoError(t, err)
		require.NotEmpty(t, offer)
		require.Equal(t, sampleMsgComment, offer.Comment)
	})

	t.Run("test propose presentation failure", func(t *testing.T) {
		oobSvc := &mockoutofband.MockOobService{
			AcceptInvitationHandle: func(*outofbandSvc.Invitation, outofbandSvc.Options) (string, error) {
				return "", fmt.Errorf(sampleClientErr)
			},
		}
		mockctx.ServiceMap[outofbandSvc.Name] = oobSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		msg, err := vcWallet.ProposeCredential(&wallet.GenericInvitation{})
		require.Error(t, err)
		require.Empty(t, msg)
	})

	t.Run("test propose presentation failure - auth error", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		msg, err := vcWallet.ProposeCredential(&wallet.GenericInvitation{})
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, msg)
	})
}

func TestClient_RequestCredential(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("test present proof success", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		response, err := vcWallet.RequestCredential(uuid.New().String(), wallet.FromPresentation(&verifiable.Presentation{}))
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, model.AckStatusPENDING, response.Status)
	})

	t.Run("test present proof success - wait for done", func(t *testing.T) {
		thID := uuid.New().String()

		icSvc := &mockissuecredential.MockIssueCredentialSvc{
			RegisterMsgEventHandle: func(ch chan<- service.StateMsg) error {
				ch <- service.StateMsg{
					Type:    service.PostState,
					StateID: "done",
					Properties: &mockdidexchange.MockEventProperties{
						Properties: map[string]interface{}{
							webRedirectStatusKey: model.AckStatusOK,
							webRedirectURLKey:    exampleWebRedirect,
						},
					},
					Msg: &mockMsg{thID: thID},
				}

				return nil
			},
		}
		mockctx.ServiceMap[issuecredentialsvc.Name] = icSvc

		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
		require.NoError(t, err)
		defer vcWallet.Close()

		response, err := vcWallet.RequestCredential(thID, wallet.FromPresentation(&verifiable.Presentation{}),
			wallet.WaitForDone(0))
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, model.AckStatusOK, response.Status)
		require.Equal(t, exampleWebRedirect, response.RedirectURL)
	})

	t.Run("test present proof failure - auth error", func(t *testing.T) {
		vcWallet, err := New(sampleUser, mockctx)
		require.NoError(t, err)
		require.NotEmpty(t, vcWallet)

		response, err := vcWallet.RequestCredential(uuid.New().String(), wallet.FromPresentation(&verifiable.Presentation{}))
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, response)
	})
}

func TestClient_ResolveCredentialManifest(t *testing.T) {
	sampleUser := uuid.New().String()
	mockctx := newMockProvider(t)

	err := CreateProfile(sampleUser, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	vcWallet, err := New(sampleUser, mockctx)
	require.NoError(t, err)
	require.NotEmpty(t, vcWallet)

	err = vcWallet.Open(wallet.WithUnlockByPassphrase(samplePassPhrase))
	require.NoError(t, err)

	t.Run("test resolving credential response", func(t *testing.T) {
		resolved, err := vcWallet.ResolveCredentialManifest(testdata.CredentialManifestMultipleVCs,
			wallet.ResolveRawResponse(testdata.CredentialResponseWithMultipleVCs))
		require.NoError(t, err)
		require.Len(t, resolved, 2)
	})

	t.Run("test resolving credential by descriptor ID", func(t *testing.T) {
		resolved, err := vcWallet.ResolveCredentialManifest(testdata.CredentialManifestMultipleVCs,
			wallet.ResolveRawCredential("udc_output", testdata.SampleUDCVC))
		require.NoError(t, err)
		require.Len(t, resolved, 1)
	})

	t.Run("test failure while resolving (closed wallet)", func(t *testing.T) {
		require.True(t, vcWallet.Close())

		resolved, err := vcWallet.ResolveCredentialManifest(testdata.CredentialManifestMultipleVCs,
			wallet.ResolveRawCredential("udc_output", testdata.SampleUDCVC))
		require.True(t, errors.Is(err, ErrWalletLocked))
		require.Empty(t, resolved)
	})
}

func newMockProvider(t *testing.T) *mockprovider.Provider {
	t.Helper()

	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	serviceMap := map[string]interface{}{
		presentproofSvc.Name:    &mockpresentproof.MockPresentProofSvc{},
		outofbandSvc.Name:       &mockoutofband.MockOobService{},
		didexchange.DIDExchange: &mockdidexchange.MockDIDExchangeSvc{},
		mediator.Coordination:   &mockmediator.MockMediatorSvc{},
		issuecredentialsvc.Name: &mockissuecredential.MockIssueCredentialSvc{},
		oobv2.Name:              &mockoutofbandv2.MockOobService{},
	}

	return &mockprovider.Provider{
		StorageProviderValue:              mockstorage.NewMockStoreProvider(),
		ProtocolStateStorageProviderValue: mockstorage.NewMockStoreProvider(),
		DocumentLoaderValue:               loader,
		ServiceMap:                        serviceMap,
	}
}

func createSampleProfile(t *testing.T, mockctx *mockprovider.Provider) {
	t.Helper()

	err := CreateProfile(sampleUserID, mockctx, wallet.WithPassphrase(samplePassPhrase))
	require.NoError(t, err)

	vcWallet, err := New(sampleUserID, mockctx)
	require.NoError(t, err)
	require.NotEmpty(t, vcWallet)
}

type mockStorageProvider struct {
	*mockstorage.MockStoreProvider
	config  storage.StoreConfiguration
	failure error
}

func (s *mockStorageProvider) SetStoreConfig(name string, config storage.StoreConfiguration) error {
	s.config = config

	return s.failure
}

func (s *mockStorageProvider) GetStoreConfig(name string) (storage.StoreConfiguration, error) {
	return s.config, nil
}

// mockMsg containing custom parent thread ID.
type mockMsg struct {
	*service.DIDCommMsgMap
	thID string
}

func (m *mockMsg) ParentThreadID() string {
	return m.thID
}

func (m *mockMsg) ThreadID() (string, error) {
	return m.thID, nil
}
