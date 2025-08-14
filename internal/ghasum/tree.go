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

import (
	"iter"

	"github.com/chains-project/ghasum/internal/gha"
)

type tree struct {
	value    *gha.GitHubAction
	children []*tree
}

func (t *tree) add(c *tree) {
	t.children = append(t.children, c)
}

func (t *tree) All() iter.Seq[gha.GitHubAction] {
	return func(yield func(gha.GitHubAction) bool) {
		_ = t.every(yield)
	}
}

func (t *tree) every(f func(gha.GitHubAction) bool) bool {
	if t.value != nil && !f(*t.value) {
		return false
	}

	for _, child := range t.children {
		if !child.every(f) {
			return false
		}
	}

	return true
}
