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
	_ "embed"
	"os"
	"path/filepath"

	"github.com/chains-project/ghasum/internal/gha"
)

var (
	//go:embed action.yml
	manifest []byte

	//go:embed action.js
	script []byte
)

func Create(root string) error {
	actionPath := filepath.Join(root, gha.GitHubDir, "actions", "ghasum")
	if err := os.MkdirAll(actionPath, 0o750); err != nil {
		return ErrCreateDir
	}

	manifestPath := filepath.Join(actionPath, gha.ManifestYml)
	if err := os.WriteFile(manifestPath, manifest, 0o660); err != nil {
		_ = os.RemoveAll(actionPath)
		return ErrCreateManifest
	}

	preScriptPath := filepath.Join(actionPath, "pre.js")
	if err := os.WriteFile(preScriptPath, script, 0o660); err != nil {
		_ = os.RemoveAll(actionPath)
		return ErrCreateScript
	}

	mainScriptPath := filepath.Join(actionPath, "main.js")
	if err := os.WriteFile(mainScriptPath, []byte{}, 0o660); err != nil {
		_ = os.RemoveAll(actionPath)
		return ErrCreateScript
	}

	return nil
}
