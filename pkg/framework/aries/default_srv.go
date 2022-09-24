/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aries

import (
	"github.com/markcryptohash/aries-framework-go/component/storageutil/mem"
	"github.com/markcryptohash/aries-framework-go/spi/storage"
)

func storeProvider() storage.Provider {
	return mem.NewProvider()
}
