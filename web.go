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
	"fmt"
	"io"
	"net/http"
)

func handlerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var (
		buf *bytes.Buffer
		err error
	)

	switch r.URL.Path[1:] { // 1 is the '/'
	default:
		http.Error(w,
			fmt.Sprintf("the requested procedure %q is not available", r.URL.Path[1:]),
			http.StatusNotFound)
		return
	case "analyze":
		buf, err = doAnalyze(w, r)
	}
	// Other procedures can be added with a similar case statement
	// above.  The buf and err variables above will be set and
	// the below code will take care of the rest.

	if err != nil {
		if e := jsonError(w, err); e != nil {
			log.Printf("ERROR: replying error: %s", e.Error())
		}
	} else if buf != nil {
		if err := jsonData(w, buf); err != nil {
			log.Printf("ERROR: replying data: %s", err.Error())
		}
	} else {
		log.Fatalf("FATAL: unreachable: err xor buf shouldn't be nil")
	}
}

func jsonError(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusInternalServerError)
	s := struct {
		Error string `json:"error"`
	}{Error: err.Error()}
	if e := json.NewEncoder(w).Encode(&s); e != nil {
		return e
	}
	return nil
}

func jsonData(w http.ResponseWriter, buf *bytes.Buffer) error {
	if _, err := fmt.Fprintf(w, `{"data":`); err != nil {
		return err
	}
	if _, err := io.Copy(w, buf); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "}"); err != nil {
		return err
	}
	return nil
}
