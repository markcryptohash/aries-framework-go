/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package vdr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/markcryptohash/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/markcryptohash/aries-framework-go/pkg/framework/aries/api/vdr"
	mockvdr "github.com/markcryptohash/aries-framework-go/pkg/mock/vdr"
)

func TestRegistry_New(t *testing.T) {
	t.Run("test new success", func(t *testing.T) {
		registry := New()
		require.NotNil(t, registry)
	})
	t.Run("test new with opts success", func(t *testing.T) {
		const sampleSvcType = "sample-svc-type"
		const sampleSvcEndpoint = "sample-svc-endpoint"
		registry := New(WithDefaultServiceEndpoint(sampleSvcEndpoint), WithDefaultServiceType(sampleSvcType))
		require.NotNil(t, registry)
		require.Equal(t, sampleSvcEndpoint, registry.defServiceEndpoint)
		require.Equal(t, sampleSvcType, registry.defServiceType)
	})
}

func TestRegistry_Close(t *testing.T) {
	t.Run("test success", func(t *testing.T) {
		registry := New()
		require.NoError(t, registry.Close())
	})
	t.Run("test error", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{CloseErr: fmt.Errorf("close error")}))
		err := registry.Close()
		require.Error(t, err)
		require.Contains(t, err.Error(), "close error")
	})
}

func TestRegistry_Resolve(t *testing.T) {
	t.Run("test invalid did input", func(t *testing.T) {
		registry := New()
		d, err := registry.Resolve("id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "wrong format did input")
		require.Nil(t, d)
	})

	t.Run("test did method not supported", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: false}))
		d, err := registry.Resolve("1:id:123")
		require.Error(t, err)
		require.Contains(t, err.Error(), "did method id not supported for vdr")
		require.Nil(t, d)
	})

	t.Run("test DID not found", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, ReadFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
				return nil, vdrapi.ErrNotFound
			},
		}))
		d, err := registry.Resolve("1:id:123")
		require.Error(t, err)
		require.Contains(t, err.Error(), vdrapi.ErrNotFound.Error())
		require.Nil(t, d)
	})

	t.Run("test error from resolve did", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, ReadFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
				return nil, fmt.Errorf("read error")
			},
		}))
		d, err := registry.Resolve("1:id:123")
		require.Error(t, err)
		require.Contains(t, err.Error(), "read error")
		require.Nil(t, d)
	})

	t.Run("test opts passed", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, ReadFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
				didOpts := &vdrapi.DIDMethodOpts{Values: make(map[string]interface{})}
				// Apply options
				for _, opt := range opts {
					opt(didOpts)
				}

				require.NotNil(t, didOpts.Values["k1"])
				return nil, nil
			},
		}))
		_, err := registry.Resolve("1:id:123", vdrapi.WithOption("k1", "v1"))
		require.NoError(t, err)
	})

	t.Run("test success", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: true}))
		_, err := registry.Resolve("1:id:123")
		require.NoError(t, err)
	})
}

func TestRegistry_Update(t *testing.T) {
	t.Run("test invalid did input", func(t *testing.T) {
		registry := New()
		err := registry.Update(&did.Doc{ID: "id"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "wrong format did input")
	})

	t.Run("test did method not supported", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: false}))
		err := registry.Update(&did.Doc{ID: "1:id:123"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "did method id not supported for vdr")
	})

	t.Run("test error from update did", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, UpdateFunc: func(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) error {
				return fmt.Errorf("update error")
			},
		}))
		err := registry.Update(&did.Doc{ID: "1:id:123"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "update error")
	})

	t.Run("test opts passed", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, UpdateFunc: func(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) error {
				didOpts := &vdrapi.DIDMethodOpts{Values: make(map[string]interface{})}
				// Apply options
				for _, opt := range opts {
					opt(didOpts)
				}

				require.NotNil(t, didOpts.Values["k1"])
				return nil
			},
		}))

		err := registry.Update(&did.Doc{ID: "1:id:123"}, vdrapi.WithOption("k1", "v1"))
		require.NoError(t, err)
	})

	t.Run("test success", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: true}))
		err := registry.Update(&did.Doc{ID: "1:id:123"})
		require.NoError(t, err)
	})
}

func TestRegistry_Deactivate(t *testing.T) {
	t.Run("test invalid did input", func(t *testing.T) {
		registry := New()
		err := registry.Deactivate("id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "wrong format did input")
	})

	t.Run("test did method not supported", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: false}))
		err := registry.Deactivate("1:id:123")
		require.Error(t, err)
		require.Contains(t, err.Error(), "did method id not supported for vdr")
	})

	t.Run("test error from deactivate did", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, DeactivateFunc: func(didID string, opts ...vdrapi.DIDMethodOption) error {
				return fmt.Errorf("deactivate error")
			},
		}))
		err := registry.Deactivate("1:id:123")
		require.Error(t, err)
		require.Contains(t, err.Error(), "deactivate error")
	})

	t.Run("test opts passed", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true, DeactivateFunc: func(didID string, opts ...vdrapi.DIDMethodOption) error {
				didOpts := &vdrapi.DIDMethodOpts{Values: make(map[string]interface{})}
				// Apply options
				for _, opt := range opts {
					opt(didOpts)
				}

				require.NotNil(t, didOpts.Values["k1"])
				return nil
			},
		}))

		err := registry.Deactivate("1:id:123", vdrapi.WithOption("k1", "v1"))
		require.NoError(t, err)
	})

	t.Run("test success", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: true}))
		err := registry.Deactivate("1:id:123")
		require.NoError(t, err)
	})
}

func TestRegistry_Create(t *testing.T) {
	t.Run("test did method not supported", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{AcceptValue: false}))
		d, err := registry.Create("id", &did.Doc{ID: "did"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "did method id not supported for vdr")
		require.Nil(t, d)
	})
	t.Run("test opts is passed", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true,
			CreateFunc: func(didDoc *did.Doc,
				opts ...vdrapi.DIDMethodOption) (doc *did.DocResolution, e error) {
				require.Equal(t, "key1", didDoc.VerificationMethod[0].ID)
				return &did.DocResolution{DIDDocument: &did.Doc{ID: "1:id:123"}}, nil
			},
		}))
		_, err := registry.Create("id", &did.Doc{VerificationMethod: []did.VerificationMethod{{ID: "key1"}}})
		require.NoError(t, err)
	})
	t.Run("with KMS opts - test opts is passed ", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true,
			CreateFunc: func(didDoc *did.Doc,
				opts ...vdrapi.DIDMethodOption) (doc *did.DocResolution, e error) {
				require.Equal(t, "key1", didDoc.VerificationMethod[0].ID)
				return &did.DocResolution{DIDDocument: &did.Doc{ID: "1:id:123"}}, nil
			},
		}))
		_, err := registry.Create("id", &did.Doc{VerificationMethod: []did.VerificationMethod{{ID: "key1"}}})
		require.NoError(t, err)
	})
	t.Run("test error from build doc", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true,
			CreateFunc: func(didDoc *did.Doc,
				opts ...vdrapi.DIDMethodOption) (doc *did.DocResolution, e error) {
				return nil, fmt.Errorf("build did error")
			},
		}))
		d, err := registry.Create("id", &did.Doc{ID: "did"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "build did error")
		require.Nil(t, d)
	})
	t.Run("test success", func(t *testing.T) {
		registry := New(WithVDR(&mockvdr.MockVDR{
			AcceptValue: true,
			CreateFunc: func(didDoc *did.Doc,
				opts ...vdrapi.DIDMethodOption) (doc *did.DocResolution, e error) {
				return &did.DocResolution{DIDDocument: &did.Doc{ID: "1:id:123"}}, nil
			},
		}))
		_, err := registry.Create("id", &did.Doc{ID: "did"})
		require.NoError(t, err)
	})
}
