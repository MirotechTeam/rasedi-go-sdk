package rasedi

type CreatePaymentPayload struct {
	Amount                     string    `json:"amount"`
	Gateways                   []Gateway `json:"gateways"`
	Title                      string    `json:"title"`
	Description                string    `json:"description"`
	RedirectURL                string    `json:"redirectUrl"`
	CallbackURL                string    `json:"callbackUrl"`
	CollectFeeFromCustomer     bool      `json:"collectFeeFromCustomer"`
	CollectCustomerEmail       bool      `json:"collectCustomerEmail"`
	CollectCustomerPhoneNumber bool      `json:"collectCustomerPhoneNumber"`
}

type CreatePaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

type CreatePaymentResponse struct {
	Body       CreatePaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

type PaymentDetailsResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

type PaymentDetailsResponse struct {
	Body       PaymentDetailsResponseBody `json:"body"`
	Headers    map[string]string          `json:"headers"`
	StatusCode int                        `json:"statusCode"`
}

type CancelPaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

type CancelPaymentResponse struct {
	Body       CancelPaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

type VerifyPayload struct {
	KeyID   string `json:"keyId"`
	Content *string `json:"content"`
}

type VerifyPaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

type VerifyPaymentResponse struct {
	Body       VerifyPaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

type PublicKey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type PublicKeysResponse struct {
	Body       []PublicKey       `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}
