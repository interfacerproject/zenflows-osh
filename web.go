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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/clone" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m := &struct{ Repo *string }{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(m); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, "Parsing the body failed: %s\n", err.Error())
		return
	}

	if m.Repo == nil || *m.Repo == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, "Please provide a repository URL to clone\n")
		return
	}

	cloneAndAnalyze(w, *m.Repo)
}

func cloneAndAnalyze(w http.ResponseWriter, repoURL string) {
	log.Printf("Cloning and Analyzing repository %s\n", repoURL)

	tmpDir, err := os.MkdirTemp("", "repo-")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating temporary directory: %s\n", err.Error())
		return
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error cloning repository: %s\n", err.Error())
		return
	}

	cmd = exec.Command("./osh", "-fC", tmpDir, "check", "--report-json=/dev/stdout")
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error running osh tool: %s\n", err.Error())
		return
	}

	output = []byte(strings.TrimPrefix(string(output), "JObject"))
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(output); err != nil {
		log.Printf("Writing output failed: %s\n", err.Error())
		return
	}

	log.Printf("Cloned and Analyzed repository %s\n", repoURL)
}
