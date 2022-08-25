package services

import (
	"github.com/btcsuite/btcutil"
	"github.com/patogallaiov/boltz-poc/connectors"
	"github.com/patogallaiov/boltz-poc/storage"
)

type Invoice struct {
	Lninvoice string `json:"lninvoice"`
}

type CheckoutService struct {
	boltz connectors.BoltzConnector
	db    storage.DBConnector
}

func NewCheckoutService(boltz connectors.BoltzConnector, db storage.DBConnector) *CheckoutService {
	return &CheckoutService{
		boltz,
		db,
	}
}

func (service *CheckoutService) CreateInvoice(request storage.PaymentRequest) (*Invoice, error) {
	response, err := service.boltz.NewReverseSwap("BTC/rBTC", "sell", btcutil.Amount(request.Amount), "", nil)
	if err != nil {
		return &Invoice{}, err
	}
	service.db.SavePayment(&storage.Payment{
		PreimageHash: response.PreimageHash,
		Preimage:     response.Preimage,
		Invoice:      response.Invoice,
		Amount:       request.Amount,
	})
	return &Invoice{
		Lninvoice: response.Invoice,
	}, nil
}

func (service *CheckoutService) GetInvoices() (any, error) {
	return service.db.GetPayments()
}

func (service *CheckoutService) GetInvoice(preimageHash string) (any, error) {
	return service.db.GetPayment(preimageHash)
}
