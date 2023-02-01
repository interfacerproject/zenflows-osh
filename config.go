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
	"net"
	"os"
)

type config struct {
	addr string
}

func loadConfig() (*config, error) {
	var c config

	addrStr, ok := os.LookupEnv("ADDR")
	if ok {
		host, port, err := net.SplitHostPort(addrStr)
		if err != nil {
			return nil, err
		}
		c.addr = net.JoinHostPort(host, port)
	} else {
		log.Printf(`WARNING: ADDR is not set; defaulting to ":7000", which will bind to all IP addresses`)
		c.addr = ":7000"
	}

	return &c, nil
}
