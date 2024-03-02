// Copyright 2024 Eric Cornelissen
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
	"errors"
	"io"
	"io/fs"

	"github.com/ericcornelissen/ghasum/internal/cache"
	"github.com/ericcornelissen/ghasum/internal/checksum"
	"github.com/ericcornelissen/ghasum/internal/sumfile"
)

type (
	// Config is the configuration for a ghasum operation.
	Config struct {
		// Repo is a pointer to the file system hierarchy of the target repository
		// for the operation.
		Repo fs.FS

		// Path is the absolute or relate path to the target repository for the
		// operation.
		//
		// This must be provided in addition to Repo because that does not allow for
		// non-read file system operation.
		Path string

		// Workflow is the file path (relative to Path) of the workflow that is the
		// subject of the operation. If this has the zero value all workflows in the
		// Repo will collectively be the subject of the operation instead.
		Workflow string

		// Cache is the cache that should be used for the operation.
		Cache cache.Cache
	}

	// Problem represents an issue detected when verifying ghasum checksums.
	Problem string
)

// Initialize will initialize ghasum for the repository specified in the given
// configuration.
func Initialize(cfg *Config) error {
	file, err := create(cfg.Path)
	if err != nil {
		return err
	}

	defer func() {
		deinitialize := (err != nil)
		if err = file.Close(); err != nil || deinitialize {
			_ = remove(cfg.Path)
		}
	}()

	actions, err := find(cfg)
	if err != nil {
		return err
	}

	checksums, err := compute(cfg, actions, checksum.BestAlgo)
	if err != nil {
		return err
	}

	content, err := encode(sumfile.VersionLatest, checksums)
	if err != nil {
		return err
	}

	if err := write(file, content); err != nil {
		return err
	}

	if err := unlock(cfg.Path); err != nil {
		return err
	}

	return nil
}

// Update will update the ghasum checksums for the repository specified in the
// given configuration.
func Update(cfg *Config, force bool) error {
	file, err := open(cfg.Path)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	raw, err := io.ReadAll(file)
	if err != nil {
		return errors.Join(ErrSumfileRead, err)
	}

	version, err := version(raw)
	if err != nil {
		if !force {
			return errors.Join(ErrSumfileRead, err)
		}

		if errors.Is(err, sumfile.ErrHeaders) || errors.Is(err, sumfile.ErrVersion) {
			version = sumfile.VersionLatest
		}
	}

	actions, err := find(cfg)
	if err != nil {
		return err
	}

	checksums, err := compute(cfg, actions, checksum.BestAlgo)
	if err != nil {
		return err
	}

	encoded, err := encode(version, checksums)
	if err != nil {
		return err
	}

	if err := clear(file); err != nil {
		return err
	}

	if err := write(file, encoded); err != nil {
		return err
	}

	if err := unlock(cfg.Path); err != nil {
		return err
	}

	return nil
}

// Verify will compare the stored ghasum checksums against recomputed checksums
// for the repository specified in the given configuration.
//
// Verification report checksums that do not match and checksums that are
// missing. It does not report checksums that are not used.
func Verify(cfg *Config) ([]Problem, error) {
	raw, err := read(cfg.Repo)
	if err != nil {
		return nil, err
	}

	stored, err := decode(raw)
	if err != nil {
		return nil, err
	}

	actions, err := find(cfg)
	if err != nil {
		return nil, err
	}

	fresh, err := compute(cfg, actions, checksum.Sha256)
	if err != nil {
		return nil, err
	}

	result := compare(fresh, stored)
	return result, nil
}
