package rasedi

// CreatePaymentPayload contains the data required to initialize a payment.
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

// CreatePaymentResponseBody represents the response body from a payment creation request.
type CreatePaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

// CreatePaymentResponse wraps the API response for creating a payment, including headers and status code.
type CreatePaymentResponse struct {
	Body       CreatePaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

// PaymentDetailsResponseBody represents the details of a specific payment.
type PaymentDetailsResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

// PaymentDetailsResponse wraps the API response for fetching payment details.
type PaymentDetailsResponse struct {
	Body       PaymentDetailsResponseBody `json:"body"`
	Headers    map[string]string          `json:"headers"`
	StatusCode int                        `json:"statusCode"`
}

// CancelPaymentResponseBody represents the response body when cancelling a payment.
type CancelPaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Amount        string        `json:"amount"`
	PaidVia       *string       `json:"paidVia"`
	PaidAt        *string       `json:"paidAt"`
	RedirectURL   string        `json:"redirectUrl"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

// CancelPaymentResponse wraps the API response for cancelling a payment.
type CancelPaymentResponse struct {
	Body       CancelPaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

// VerifyPayload contains the data required to verify a webhook signature or payment completion.
type VerifyPayload struct {
	KeyID   string  `json:"keyId"`
	Content *string `json:"content"`
}

// VerifyPaymentResponseBody represents the verification status of a payment.
type VerifyPaymentResponseBody struct {
	ReferenceCode string        `json:"referenceCode"`
	Status        PaymentStatus `json:"status"`
	PayoutAmount  *string       `json:"payoutAmount"`
}

// VerifyPaymentResponse wraps the API response for payment verification.
type VerifyPaymentResponse struct {
	Body       VerifyPaymentResponseBody `json:"body"`
	Headers    map[string]string         `json:"headers"`
	StatusCode int                       `json:"statusCode"`
}

// PublicKey represents a public key used for signature verification.
type PublicKey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// PublicKeysResponse wraps the API response for fetching public keys.
type PublicKeysResponse struct {
	Body       []PublicKey       `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}
