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
	"encoding/json"
	"fmt"
	"net/http"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("access-control-allow-origin", "*")
		w.Header().Add("access-control-allow-credentials", "true")
		w.Header().Add("access-control-allow-methods", http.MethodPost)
		w.Header().Add("access-control-allow-headers", "content-type, content-length, accept")

		if r.Method == http.MethodOptions {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func handlerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var (
		buf string
		err error
	)

	switch r.URL.Path[1:] { // 1 is the '/'
	default:
		err = fmt.Errorf("the requested procedure %q is not available", r.URL.Path[1:])
		if e := jsonErr(w, err, http.StatusNotFound); e != nil {
			log.Printf("ERROR: jsonErr(): %s", e)
		}
		return
	case "analyze":
		buf, err = doAnalyze(w, r)
	}
	// Other procedures can be added with a similar case statement
	// above.  The buf and err variables above will be set and
	// the below code will take care of the rest.

	if err != nil {
		if e := jsonErr(w, err, http.StatusInternalServerError); e != nil {
			log.Printf("ERROR: jsonErr(): %s", e)
		}
	} else if buf != "" {
		if err := jsonOk(w, buf); err != nil {
			log.Printf("ERROR: jsonOk(): %s", err)
		}
	} else {
		log.Fatalf("FATAL: unreachable: err xor buf shouldn't be nil")
	}
}

func jsonErr(w http.ResponseWriter, err error, stat int) error {
	w.WriteHeader(stat)
	s := struct {
		Err string `json:"err"`
	}{err.Error()}
	if e := json.NewEncoder(w).Encode(&s); e != nil {
		return e
	}
	return nil
}

func jsonOk(w http.ResponseWriter, s string) error {
	if _, err := fmt.Fprintf(w, `{"ok":%s}`, s); err != nil {
		return err
	}

	return nil
}
