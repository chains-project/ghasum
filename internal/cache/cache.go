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

package cache

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Cache represents a cache located on the file system.
type Cache struct {
	// Path is the path of the cache on the file system.
	path string

	// Ephemeral marks the cache as such, locating it in the system's temporary
	// directory.
	ephemeral bool
}

// Cleanup the cache if it is ephemeral, removing it.
//
// Any errors are silently ignored under the assumption that the temporary
// directory will be cleaned automatically.
func (c *Cache) Cleanup() {
	if c.ephemeral {
		_ = c.Clear()
	}
}

// Clear the contents of the cache, removing everything.
func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.path); err != nil {
		return fmt.Errorf("could not clear %q: %v", c.path, err)
	}

	return nil
}

// Evict old entries from the cache, removing them.
func (c *Cache) Evict() (uint, error) {
	deadline := time.Now().AddDate(0, 0, -5)

	var count uint
	walk := func(path string, entry fs.DirEntry, _ error) error {
		depth := strings.Count(path, string(os.PathSeparator))
		if depth < 2 {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("could not get file info for %q", path)
		}

		if info.ModTime().Before(deadline) {
			count += 1
			_ = os.RemoveAll(filepath.Join(c.path, path))
			return fs.SkipDir
		}

		return fs.SkipDir
	}

	fsys, err := os.OpenRoot(c.path)
	if errors.Is(err, fs.ErrNotExist) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("could not open cache directory: %v", err)
	}

	if err := fs.WalkDir(fsys.FS(), ".", walk); err != nil {
		return 0, fmt.Errorf("cache eviction failed: %v", err)
	}

	return count, nil
}

// Init the cache.
func (c *Cache) Init() error {
	if c.ephemeral {
		location, err := os.MkdirTemp(os.TempDir(), "ghasum-clone-*")
		if err != nil {
			return fmt.Errorf("could not create temporary cache: %v", err)
		}

		c.path = location
	} else {
		if err := os.MkdirAll(c.path, 0o700); err != nil {
			return fmt.Errorf("could not create cache at %q: %v", c.path, err)
		}
	}

	return nil
}

// Path returns the path to the cache on the file system.
func (c *Cache) Path() string {
	return c.path
}

// New creates an uninitialized cache.
//
// If location is an empty string the location will default to the user's home
// directory.
//
// If ephemeral is set the cache will be located in a unique directory in the
// system's temporary directory (and the given location is ignored).
func New(options ...Option) (Cache, error) {
	opts := defaultOpts
	for _, option := range options {
		opts = option(opts)
	}

	var c Cache
	switch {
	case opts.Ephemeral:
		c.ephemeral = true
	case opts.Location == "":
		home, err := os.UserHomeDir()
		if err != nil {
			return c, fmt.Errorf("could not get home directory: %v", err)
		}

		c.path = filepath.Join(home, ".ghasum")
	default:
		c.path = opts.Location
	}

	if opts.Evict {
		if _, err := c.Evict(); err != nil {
			return c, err
		}
	}

	return c, nil
}
