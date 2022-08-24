package connectors

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

const (
	retries    int = 3
	rpcSleep       = 2 * time.Second
	rpcTimeout     = 5 * time.Second
	ethSleep       = 5 * time.Second
	ethTimeout     = 5 * time.Minute

	newAccountGasCost = uint64(25000)
)

type RSKConnector interface {
	Connect(endpoint string, chainId *big.Int) error
	CheckConnection() error
	Subscribe([]common.Address, func(gethTypes.Log)) error
	Close()
	GetChainId() (*big.Int, error)
	EstimateGas(addr string, value *big.Int, data []byte) (uint64, error)
	GasPrice() (*big.Int, error)
	GetTxStatus(ctx context.Context, tx *gethTypes.Transaction) (bool, error)
}

type RSK struct {
	c          *ethclient.Client
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func NewRSK(privateKeyString string) (*RSK, error) {
	key, _ := crypto.HexToECDSA(privateKeyString)
	address := crypto.PubkeyToAddress(key.PublicKey)
	return &RSK{
		privateKey: key,
		address:    address,
	}, nil
}

func (rsk *RSK) Connect(endpoint string, chainId *big.Int) error {
	log.Debug("connecting to RSK node on ", endpoint)

	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	var ethC *ethclient.Client
	switch u.Scheme {
	case "http", "https":
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.DisableKeepAlives = true

		httpC := new(http.Client)
		httpC.Transport = transport

		c, err := rpc.DialHTTPWithClient(endpoint, httpC)
		if err != nil {
			return err
		}

		ethC = ethclient.NewClient(c)
	default:
		ethC, err = ethclient.Dial(endpoint)
		if err != nil {
			return err
		}
	}

	rsk.c = ethC

	log.Debug("verifying connection to RSK node")
	// test connection
	rskChainId, err := rsk.GetChainId()
	if err != nil {
		return err
	}
	//check chain id
	if chainId.Cmp(rskChainId) != 0 {
		return fmt.Errorf("chain id mismatch; expected chain id: %v, rsk node chain id: %v", chainId, rskChainId)
	}
	return nil
}

func (rsk *RSK) Subscribe(addresses []common.Address, callback func(log gethTypes.Log)) error {

	query := ethereum.FilterQuery{
		Addresses: addresses,
	}

	//Log Channel
	logs := make(chan types.Log)

	//Subscribe
	sub, err := rsk.c.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	// Leave a thread pooling
	go poolSubscription(sub, logs, callback)

	return nil
}

func poolSubscription(sub ethereum.Subscription, logs chan gethTypes.Log, callback func(log gethTypes.Log)) {
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error in subscription %v.", err)
		case vLog := <-logs:
			if callback == nil {
				log.Debug("Event received:")
				log.Debugln(vLog) // pointer to event log
				log.Error("Event received but no callback declared.")
			} else {
				go callback(vLog)
			}
		}
	}
}

func (rsk *RSK) CheckConnection() error {
	_, err := rsk.GetChainId()
	return err
}

func (rsk *RSK) Close() {
	log.Debug("closing RSK connection")
	rsk.c.Close()
}

func (rsk *RSK) GetChainId() (*big.Int, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var chainId *big.Int
		chainId, err = rsk.c.ChainID(ctx)
		if err == nil {
			return chainId, nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error retrieving chain id: %v", err)
}

func (rsk *RSK) EstimateGas(addr string, value *big.Int, data []byte) (uint64, error) {
	if !common.IsHexAddress(addr) {
		return 0, fmt.Errorf("invalid address: %v", addr)
	}

	dst := common.HexToAddress(addr)

	var additionalGas uint64
	if rsk.isNewAccount(dst) {
		additionalGas = newAccountGasCost
	}

	msg := ethereum.CallMsg{
		To:    &dst,
		Data:  data,
		Value: new(big.Int).Set(value),
	}

	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var gas uint64
		gas, err = rsk.c.EstimateGas(ctx, msg)
		if gas > 0 {
			return gas + additionalGas, nil
		}
		time.Sleep(rpcSleep)
	}
	return 0, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) GasPrice() (*big.Int, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var price *big.Int
		price, err = rsk.c.SuggestGasPrice(ctx)
		if price != nil && price.Cmp(big.NewInt(0)) >= 0 {
			return price, nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) GetTxStatus(ctx context.Context, tx *gethTypes.Transaction) (bool, error) {
	ticker := time.NewTicker(ethSleep)

	for {
		select {
		case <-ticker.C:
			cctx, cancel := context.WithTimeout(ctx, rpcTimeout)
			defer cancel()
			r, _ := rsk.c.TransactionReceipt(cctx, tx.Hash())
			if r != nil {
				return r.Status == 1, nil
			}
		case <-ctx.Done():
			ticker.Stop()
			return false, fmt.Errorf("operation cancelled")
		}
	}
}

func (rsk *RSK) isNewAccount(addr common.Address) bool {
	var (
		err  error
		code []byte
		bal  *big.Int
		n    uint64
	)
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		code, err = rsk.c.CodeAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		bal, err = rsk.c.BalanceAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		n, err = rsk.c.NonceAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	return len(code) == 0 && bal.Cmp(common.Big0) == 0 && n == 0
}

func DecodeRSKAddress(address string) ([]byte, error) {
	trim := strings.TrimPrefix(address, "0x")
	if !common.IsHexAddress(trim) {
		return nil, fmt.Errorf("invalid address: %v", address)
	}
	return common.HexToAddress(trim).Bytes(), nil
}

func copyHex(str string, dst []byte) error {
	bts, err := parseHex(str)
	if err != nil {
		return err
	}
	copy(dst, bts[:])
	return nil
}

func parseHex(str string) ([]byte, error) {
	bts, err := hex.DecodeString(strings.TrimPrefix(str, "0x"))
	if err != nil {
		return nil, err
	}
	return bts, nil
}
