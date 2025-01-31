/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cm

import "github.com/markcryptohash/aries-framework-go/pkg/doc/presexch"

// hasAnyAlgorithmsOrProofTypes looks at the given Format object and determines if it has any algorithms or proof types
// listed.
func hasAnyAlgorithmsOrProofTypes(format presexch.Format) bool {
	if anyJWTTypeHasAlgs(format) || anyLDPTypeHasProofTypes(format) {
		return true
	}

	return false
}

func anyJWTTypeHasAlgs(format presexch.Format) bool {
	if hasAnyAlgs(format.Jwt) ||
		hasAnyAlgs(format.JwtVC) ||
		hasAnyAlgs(format.JwtVP) {
		return true
	}

	return false
}

func anyLDPTypeHasProofTypes(format presexch.Format) bool {
	if hasAnyProofTypes(format.Ldp) ||
		hasAnyProofTypes(format.LdpVC) ||
		hasAnyProofTypes(format.LdpVP) {
		return true
	}

	return false
}

func hasAnyAlgs(jwtType *presexch.JwtType) bool {
	if jwtType != nil && len(jwtType.Alg) > 0 {
		return true
	}

	return false
}

func hasAnyProofTypes(ldpType *presexch.LdpType) bool {
	if ldpType != nil && len(ldpType.ProofType) > 0 {
		return true
	}

	return false
}

func lookUpString(model map[string]interface{}, key string) (string, bool) {
	raw, ok := model[key]
	if !ok {
		return "", false
	}

	val, ok := raw.(string)

	return val, ok
}

func lookUpMap(model map[string]interface{}, key string) (map[string]interface{}, bool) {
	raw, ok := model[key]
	if !ok {
		return nil, false
	}

	val, ok := raw.(map[string]interface{})

	return val, ok
}

func lookUpArray(model map[string]interface{}, key string) ([]interface{}, bool) {
	raw, ok := model[key]
	if !ok {
		return nil, false
	}

	val, ok := raw.([]interface{})

	return val, ok
}
