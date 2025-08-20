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

package cache

// Option is a function to configure the cache.
type Option func(Options) Options

// Options for a [New] cache.
type Options struct {
	Location  string
	Ephemeral bool
	Evict     bool
}

var defaultOpts = Options{
	Ephemeral: false,
	Evict:     true,
}

// WithEphemeralCache makes the cache ephemeral (single-run use).
func WithEphemeralCache(v bool) Option {
	return func(opts Options) Options {
		opts.Ephemeral = v
		return opts
	}
}

// WithEviction enables or disabled cache evection.
func WithEviction(v bool) Option {
	return func(opts Options) Options {
		opts.Evict = v
		return opts
	}
}

// WithLocation sets the location of the cache. The location is ignored when the
// cache is ephemeral.
func WithLocation(v string) Option {
	return func(opts Options) Options {
		opts.Location = v
		return opts
	}
}
