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

package gha

import "errors"

var (
	// ErrDockerUses is the error used when a uses value is for a Docker action.
	ErrDockerUses = errors.New("uses is a Docker Hub/Container Registry action")

	// ErrInvalidUses is the error used for an invalid uses value.
	ErrInvalidUses = errors.New("invalid uses value")

	// ErrInvalidUsesRepo is the error used for a uses value with an invalid
	// repository.
	ErrInvalidUsesRepo = errors.New("invalid repository in uses")

	// ErrInvalidUsesPath is the error used for a uses value with an invalid
	// path in the repository.
	ErrInvalidUsesPath = errors.New("invalid repository path in uses")

	// ErrLocalAction is the error used when a uses value is for a local action.
	ErrLocalAction = errors.New("uses is a local action")

	// ErrDockerfileManifest is the error used when (only) a Dockerfile Action
	// manifest is present.
	ErrDockerfileManifest = errors.New("found a Dockerfile manifest")

	// ErrNoManifest is the error used when no Action manifest could be found.
	ErrNoManifest = errors.New("could not find a manifest")
)
