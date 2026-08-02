// Copyright 2026 Eric Cornelissen
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

package action

import (
	"errors"
)

var (
	// ErrCreateDir is the error used when the local action directory could not
	// be created.
	ErrCreateDir = errors.New("local action directory could not be created")

	// ErrCreateManifest is the error used when the local action's manifest
	// could not be created.
	ErrCreateManifest = errors.New("local action manifest could not be created")

	// ErrCreateScript is the error used when the local action's script could
	// not be created.
	ErrCreateScript = errors.New("local action script could not be created")
)
