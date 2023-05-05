// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
)

func doShortlog(w http.ResponseWriter, r *http.Request) (string, error) {
	in := &struct {
		Repo *string `json:"repo"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return "", err
	}

	if in.Repo == nil || *in.Repo == "" {
		return "", errors.New("the repo url must be a non-empty string")
	}

	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", *in.Repo, tmpDir)
	cmd.Env = []string{"GIT_TERMINAL_PROMPT=0"}
	if err := cmd.Run(); err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	cmd = exec.Command("git", "shortlog", "-sne", "HEAD")
	cmd.Dir = tmpDir
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	b, err := json.Marshal(buf.String())
	if err != nil {
		return "", err
	}

	return string(b), nil
}
