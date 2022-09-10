package config

import "math/big"

type Config struct {
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
		Endpoint     string
		AbiPath      string
		ERC20AbiPath string
	}
	Accounts struct {
		RSK struct {
			PrivateKey string
			Address    string
		}
	}
}
