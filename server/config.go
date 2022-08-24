package main

import "math/big"

type config struct {
	LogFile string
	Debug   bool
	ErpKeys []string

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK struct {
		Endpoint string
		ChainId  *big.Int
	}
	Boltz struct {
		Endpoint string
	}
	Accounts struct {
		RSK struct {
			PrivateKey string
			Address    string
		}
	}
}
