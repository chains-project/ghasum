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

package cache

import "errors"

var (
	// ErrClear is the error used when ghasum could not clear the cache.
	ErrClear = errors.New("could not clear cache")

	// ErrCreate is the error used when ghasum could not create the cache
	// directory.
	ErrCreate = errors.New("could not create cache")

	// ErrEvict is the error used when ghasum could not evict at least one cache
	// entry which should have been evicted.
	ErrEvict = errors.New("cache eviction failed")

	// ErrOpen is the error used when ghasum is not able to open the cache
	// directory.
	ErrOpen = errors.New("could not open the cache directory")
)
