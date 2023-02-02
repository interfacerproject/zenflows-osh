/*
Written and maintained by srfsh <info@dyne.org>.
Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
)

func doAnalyze(w http.ResponseWriter, r *http.Request) (*bytes.Buffer, error) {
	in := &struct {
		Repo *string `json:"repo"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}

	if in.Repo == nil || *in.Repo == "" {
		return nil, errors.New("the repo url must be a non-empty string")
	}

	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", *in.Repo, tmpDir)
	cmd.Env = []string{"GIT_TERMINAL_PROMPT=0"}
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	cmd = exec.Command("osh", "-fC", tmpDir, "check", "--report-json=/dev/stdout")
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// as the output starts with JObject, we descard it
	b := make([]byte, 7)
	n, err := buf.Read(b)
	if n != 7 || string(b) != "JObject" || err != nil {
		return nil, errors.New("osh-tool is acting up")
	}

	return buf, nil
}
