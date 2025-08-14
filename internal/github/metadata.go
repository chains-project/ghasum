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

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type apiRepoMetadata struct {
	Archived bool `json:"archived"`
}

const apiUrl = "https://api.github.com"

// Archived returns whether the given repository is archived on GitHub.
func Archived(repo *Repository) (bool, error) {
	metadata, err := getRepoMetadata(repo)
	if err != nil {
		return false, err
	}

	return metadata.Archived, nil
}

func getRepoMetadata(repo *Repository) (apiRepoMetadata, error) {
	var metadata apiRepoMetadata

	url := fmt.Sprintf("%s/repos/%s/%s", apiUrl, repo.Owner, repo.Project)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	if token, ok := os.LookupEnv("GH_TOKEN"); ok {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return metadata, fmt.Errorf("GET %s failed: %v", url, err)
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return metadata, fmt.Errorf("GET %s failed with status %d", url, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return metadata, fmt.Errorf("GET %s response malformed: %v", url, err)
	}

	return metadata, nil
}
