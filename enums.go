package rasedi

// Gateway represents the payment gateway provider.
type Gateway string

const (
	GatewayFib        Gateway = "FIB"
	GatewayZain       Gateway = "ZAIN"
	GatewayAsiaPay    Gateway = "ASIA_PAY"
	GatewayFastPay    Gateway = "FAST_PAY"
	GatewayNassWallet Gateway = "NASS_WALLET"
	GatewayCreditCard Gateway = "CREDIT_CARD"
)

// PaymentStatus represents the current status of a payment.
type PaymentStatus string

const (
	PaymentStatusTimedOut PaymentStatus = "TIMED_OUT"
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusCanceled PaymentStatus = "CANCELED"
	PaymentStatusFailed   PaymentStatus = "FAILED"
)
