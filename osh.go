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
	stdlog "log"
	"net/http"
	"os"
)

var log = stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime|stdlog.LUTC)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("bad config: %s", err.Error())
	}
	log.Printf("Starting service on %s", conf.addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/clone", cloneHandler)

	log.Fatal(http.ListenAndServe(conf.addr, mux))
}
