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
	"testing"
)

func TestAlgorithms(t *testing.T) {
	t.Parallel()

	t.Run("Default value", func(t *testing.T) {
		t.Parallel()

		var algo Algo
		if got, ok := hashes[algo]; ok {
			t.Errorf("Want no default algo, got %v", got)
		}
	})

	t.Run("Supported", func(t *testing.T) {
		t.Parallel()

		testCases := []Algo{
			BestAlgo,
			Sha256,
		}

		for _, algo := range testCases {
			if _, ok := hashes[algo]; !ok {
				t.Errorf("Want an algorithm for %d, got none", algo)
			}
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		t.Parallel()

		testCases := []Algo{
			-1,
			255,
		}

		for _, algo := range testCases {
			if got, ok := hashes[algo]; ok {
				t.Errorf("Want no algorithm for %d, got %v", algo, got)
			}
		}
	})
}
