package connectors

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BadRequestError string

func (e *BadRequestError) Error() string {
	return string(*e)
}

type boltzReverseSwap struct {
	ID                 string `json:"id"`
	Invoice            string `json:"invoice"`
	RefundAddress      string `json:"refundAddress"`
	LockupAddress      string `json:"lockupAddress"`
	OnchainAmount      int64  `json:"onchainAmount"`
	TimeoutBlockHeight int64  `json:"timeoutBlockHeight"`
}

type ReverseSwap struct {
	boltzReverseSwap
	Preimage     string
	PreimageHash string
	Key          string
}

type Fees struct {
	Percentage float64
	Lockup     int64
	Claim      int64
}

type ReverseSwapInfo struct {
	FeesHash string
	Max      int64
	Min      int64
	Fees     Fees
}

type PairResponse struct {
	Warnings []string `json:"warnings"`
	Pairs    map[string]struct {
		Rate   float64 `json:"rate"`
		Hash   string  `json:"hash"`
		Limits struct {
			Maximal         int64 `json:"maximal"`
			Minimal         int64 `json:"minimal"`
			MaximalZeroConf struct {
				BaseAsset  int64 `json:"baseAsset"`
				QuoteAsset int64 `json:"quoteAsset"`
			} `json:"maximalZeroConf"`
		}
		Fees struct {
			Percentage float64 `json:"percentage"`
			MinerFees  struct {
				BaseAsset struct {
					Normal  int64 `json:"normal"`
					Reverse struct {
						Lockup int64 `json:"lockup"`
						Claim  int64 `json:"claim"`
					} `json:"reverse"`
				} `json:"baseAsset"`
				QuoteAsset struct {
					Normal  int64 `json:"normal"`
					Reverse struct {
						Lockup int64 `json:"lockup"`
						Claim  int64 `json:"claim"`
					} `json:"reverse"`
				} `json:"quoteAsset"`
			} `json:"minerFees"`
		} `json:"fees"`
	} `json:"pairs"`
}

type ContractsResponse struct {
	Ethereum struct {
		Network struct {
			ChainId int8 `json:"chainId"`
		} `json:"network"`
		SwapContracts struct {
			EtherSwap string `json:"EtherSwap"`
			ERC20Swap string `json:"ERC20Swap"`
		} `json:"swapContracts"`
		Tokens struct {
			DOC string `json:"DOC"`
		} `json:"tokens"`
	}
}

type RoutingHint struct {
	HopHintsList []struct {
		NodeID                    string `json:"nodeId"`
		ChanID                    string `json:"chanId"`
		FeeBaseMsat               uint32 `json:"feeBaseMsat"`
		FeeProportionalMillionths uint32 `json:"feeProportionalMillionths"`
		CltvExpiryDelta           uint32 `json:"cltvExpiryDelta"`
	} `json:"hopHintsList"`
}

type transactionStatus struct {
	Status      string `json:"status"`
	Transaction struct {
		ID  string `json:"id"`
		Hex string `json:"hex"`
		ETA int    `json:"eta",omitempty`
	} `json:"transaction",omitempty`
}

type transactionRequest struct {
	ID string `json:"id"`
}

type lockupEvent struct {
	Amount       *big.Int
	ClaimAddress common.Address
	Timelock     *big.Int
}
