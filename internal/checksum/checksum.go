// Copyright 2024-2025 Eric Cornelissen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checksum

import (
	"fmt"

	"golang.org/x/mod/sumdb/dirhash"
)

// Algo represents a cryptographic hash algorithm.
type Algo int

const (
	_ Algo = iota

	// Sha256 identifies the SHA256 hashing algorithm.
	Sha256

	// BestAlgo identifies the best available hashing algorithm.
	BestAlgo = Sha256
)

var hashes = map[Algo]dirhash.Hash{
	Sha256: dirhash.Hash1,
}

// Compute the checksum over the directory at the given path using the specified
// cryptographic hash algorithm.
func Compute(path string, algo Algo) (string, error) {
	hash := hashes[algo]
	checksum, err := dirhash.HashDir(path, "", hash)
	if err != nil {
		return "", fmt.Errorf("could not compute checksum: %v", err)
	}

	return checksum, nil
}
