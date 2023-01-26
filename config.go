package main

import (
	"errors"
	"net"
	"os"
)

type config struct {
	addr string
}

func loadConfig() (*config, error) {
	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		return nil, errors.New("ADDR must be set")
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	return &config{
		addr: host + ":" + port,
	}, nil
}
