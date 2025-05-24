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

package sumfile

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// An Entry represents a single checksum entry in a checksum file.
type Entry struct {
	// Checksum is the checksum value for the entry.
	Checksum string

	// ID is the identifier for the entry. Can have any number of parts but must
	// not be empty.
	ID []string
}

// Decode parses the given checksum file content into Entries. This will error
// if there is a syntax error in the checksum file or if the checksum file is
// otherwise corrupted (for example multiple checksum directives for one Entry).
func Decode(stored string) ([]Entry, error) {
	_, entries, err := parseFile(stored)
	return entries, err
}

// DecodeVersion parses the given checksum file content to extract the version.
func DecodeVersion(stored string) (Version, error) {
	headers, _, parseErr := parseFile(stored)
	version, _ := extractVersion(headers)
	if parseErr != nil {
		return version, parseErr
	}

	return version, nil
}

// Encode encodes the given checksums according to the specification of the
// given version.
func Encode(version Version, checksums []Entry) (string, error) {
	var (
		encoded string
		err     error
	)

	switch version {
	case Version1:
		encoded, err = encodeV1(checksums)
	default:
		err = unknownVersion(version)
	}

	return fmt.Sprintf("version %d\n\n%s", version, encoded), err
}

func parseFile(stored string) (map[string]string, []Entry, error) {
	lines := strings.Split(stored, "\n")
	if strings.Contains(stored, "\r") {
		lines = strings.Split(stored, "\r\n")
	}

	headers, err := parseHeaders(lines)
	if err != nil {
		return nil, nil, err
	}

	version, err := extractVersion(headers)
	if err != nil {
		return headers, nil, err
	}

	if lines[len(lines)-1] != "" {
		err = errors.New("missing final newline")
		return headers, nil, errors.Join(ErrSyntax, err)
	}

	content := []string{}
	if len(lines) > len(headers)+1 {
		content = lines[len(headers)+1 : len(lines)-1]
	}

	var entries []Entry
	switch version {
	case Version1:
		entries, err = decodeV1(content)
	default:
		err = unknownVersion(version)
	}

	if err != nil {
		return headers, entries, err
	}

	return headers, entries, nil
}

func parseHeaders(lines []string) (map[string]string, error) {
	headers := make(map[string]string, 0)
	for i, line := range lines {
		if len(line) == 0 {
			break
		}

		j := strings.IndexRune(line, ' ')
		if j == -1 {
			err := fmt.Errorf("invalid header on line %d", i+1)
			return nil, errors.Join(ErrHeaders, err)
		}

		key := line[0:j]
		value := line[j+1:]
		if _, ok := headers[key]; ok {
			err := fmt.Errorf("duplicate header %q on line %d", key, i+1)
			return nil, errors.Join(ErrHeaders, err)
		}

		headers[key] = value
	}

	return headers, nil
}

func extractVersion(headers map[string]string) (Version, error) {
	version, ok := headers["version"]
	if !ok {
		err := errors.New("version not found")
		return 0, errors.Join(ErrVersion, err)
	}

	rawVersion, err := strconv.Atoi(version)
	if err != nil {
		err := errors.New("version not a number")
		return 0, errors.Join(ErrVersion, err)
	}

	return Version(rawVersion), nil
}

func unknownVersion(version Version) error {
	err := fmt.Errorf("unknown version %d", version)
	return errors.Join(ErrVersion, err)
}
