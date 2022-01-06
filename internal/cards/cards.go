package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// トランザクションに情報を追加する場合はこのようにする
	params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

func cardErrMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "カードが拒否されました"
	case stripe.ErrorCodeExpiredCard:
		msg = "カードの有効期限が切れています"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "カードのCVCが間違っています"
	case stripe.ErrorCodeIncorrectZip:
		msg = "カードの郵便番号が間違っています"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "支払い金額が大きすぎます"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "残高が不足しています"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "郵便番号が間違っています"
	default:
		msg = "カードが拒否されました"

	}

	return msg
}
