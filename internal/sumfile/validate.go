// Copyright 2024-2026 Eric Cornelissen
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

package sumfile

import (
	"strings"
)

func hasMissing(entries []Entry) bool {
	for _, entry := range entries {
		if len(entry.Checksum) == 0 || len(entry.ID) == 0 {
			return true
		}

		for _, part := range entry.ID {
			if len(part) == 0 {
				return true
			}
		}
	}

	return false
}

func hasDuplicates(entries []Entry) bool {
	seen := make(map[string]any, 0)
	for _, entry := range entries {
		var tmp strings.Builder
		for _, part := range entry.ID {
			tmp.WriteString(part)
		}

		key := tmp.String()
		if _, ok := seen[key]; ok {
			return true
		}

		seen[key] = struct{}{}
	}

	return false
}
