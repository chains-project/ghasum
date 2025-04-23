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
	ErrDockerUses      = errors.New("uses is a Docker Hub/Container Registry action")
	ErrInvalidUses     = errors.New("invalid uses value")
	ErrInvalidUsesRepo = errors.New("invalid repository in uses")
	ErrInvalidUsesPath = errors.New("invalid repository path in uses")
	ErrLocalAction     = errors.New("uses is a local action")
)
