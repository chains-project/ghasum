// Copyright 2025 Eric Cornelissen
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

package ghasum

// UpdateReport is a report produced by [Update].
type UpdateReport struct {
	// The number of checksums that were added.
	Added uint

	// The number of checksums that were not changed
	Kept uint

	// the number of checksums that were override (forced updates).
	Overridden uint

	// The number of checksums that were removed.
	Removed uint

	// The number of checksums that were updated.
	Updated uint
}

// VerifyReport is a report produced by [Verify].
type VerifyReport struct {
	// The list of problems that occurred during verification.
	Problems []Problem

	// The total number of actions that were verified.
	Total int
}
