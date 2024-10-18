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

import "errors"

var (
	// ErrInitialized is the error used when ghasum is not expected to be
	// initialized but is.
	ErrInitialized = errors.New("ghasum is already initialized")

	// ErrNotInitialized is the error used when ghasum is expected to be
	// initialized but is not.
	ErrNotInitialized = errors.New("ghasum has not yet been initialized")

	// ErrSumfileCreate is the error used when the ghasum checksum file could not
	// be created.
	ErrSumfileCreate = errors.New("could not create a checksum file")

	// ErrSumfileEncode is the error used when the ghasum checksum file could not
	// be encoded.
	ErrSumfileEncode = errors.New("could not encode the checksum file")

	// ErrSumfileOpen is the error used when the ghasum checksum file could not
	// be opened.
	ErrSumfileOpen = errors.New("could not open the checksum file")

	// ErrSumfileDecode is the error used when a ghasum checksum file could not
	// be decoded.
	ErrSumfileDecode = errors.New("could not decode the checksum file")

	// ErrSumfileRead is the error used when the ghasum checksum file could not
	// be read.
	ErrSumfileRead = errors.New("could not read from the checksum file")

	// ErrSumfileRemove is the error used when the ghasum checksum file could not
	// be removed.
	ErrSumfileRemove = errors.New("could not remove the checksum file")

	// ErrSumfileUnlock is the error used when the ghasum checksum file could not
	// be unlocked after usage.
	ErrSumfileUnlock = errors.New("could not unlock the checksum file")

	// ErrSumfileWrite is the error used when the ghasum checksum file could not
	// be written to.
	ErrSumfileWrite = errors.New("could not write to the checksum file")
)
