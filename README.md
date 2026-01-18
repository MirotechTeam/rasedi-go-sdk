# Rasedi Go SDK

The official Go SDK for [Rasedi Payment Gateway](https://rasedi.com). This library provides a simple and powerful way to interact with the Rasedi API to create, verify, and manage payments in your Go applications.

## Features

- **Simple Integration**: Easy-to-use methods for all payment operations.
- **Secure**: Built-in authentication and request signing using ECDSA/RSA/Ed25519.
- **Robust Verification**: Utilities to verify JWS signatures from payment callbacks.
- **Type-Safe**: Fully typed structs for requests and responses.

## Installation

Install the package using `go get`:

```bash
go get github.com/MirotechTeam/rasedi-go-sdk
```

## Quick Start

### 1. Initialize the Client

Initialize the client with your **Private Key** and **Secret Key**.

```go
package main

import (
	"fmt"
	"log"

	rasedi "github.com/MirotechTeam/rasedi-go-sdk"
)

func main() {
	// Your credentials
	secretKey := "live_..."
	privateKey := `-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----`

	// Initialize client
	client := rasedi.NewClient(privateKey, secretKey)
	
	fmt.Println("Client initialized successfully")
}
```

### 2. Create a Payment

```go
payload := rasedi.CreatePaymentPayload{
    Amount:                 "10000", // Amount in smallest currency unit (e.g. fils or cents)
    Title:                  "Order #1234",
    Description:            "Payment for digital goods",
    Gateways:               []rasedi.Gateway{rasedi.GatewayCreditCard},
    RedirectURL:            "https://your-site.com/callback",
    CallbackURL:            "https://your-site.com/webhook",
    CollectFeeFromCustomer: true,
    CollectCustomerEmail:   true,
}

response, err := client.CreatePayment(payload)
if err != nil {
    log.Fatalf("Failed to create payment: %v", err)
}

fmt.Printf("Payment created! Redirect user to: %s\n", response.Body.RedirectURL)
fmt.Printf("Reference Code: %s\n", response.Body.ReferenceCode)
```

### 3. Verify a Payment (Webhook/Callback)

When Rasedi redirects the user back to your site or sends a webhook, you will receive a payload with a `keyId` and `content`. Use the `Verify` method to validate the signature and decode the payment details.

```go
verifyPayload := rasedi.VerifyPayload{
    KeyID:   "key_id_from_request",
    Content: &contentFromRequest, // JWS token string
}

result, err := client.Verify(verifyPayload)
if err != nil {
    log.Fatalf("Verification failed: %v", err)
}

fmt.Printf("Payment Verified: %+v\n", result["body"])
// Check status in result["body"]["status"]
```

### 4. Check Payment Status

You can also check the status of a payment manually using its reference code.

```go
details, err := client.GetPaymentByReferenceCode("reference-code-here")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Current Status: %s\n", details.Body.Status)
```

### 5. Cancel a Payment

```go
cancelResponse, err := client.CancelPayment("reference-code-here")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Payment cancelled successfully")
```

## Configuration

### Test vs Live Mode
The SDK automatically determines the environment (Test vs Live) based on your **Secret Key**. 
- If the secret key contains `"test"`, it switches to the Test environment.
- Otherwise, it defaults to Live.



## License

This project is licensed under the MIT License.
