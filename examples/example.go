package main

import (
	"encoding/json"
	"fmt"
	"log"

	rasedi "github.com/MirotechTeam/rasedi-go-sdk"
)

func main() {
	secretKey := "live_laisxVjnNnoY1w5mwWP6YwzfPg_zmu2BnWnJH1uCOzOGcAflAYShdjVPuDAG10DLSEpTOlsOopiyTJHJjO4fbqqU"
	privateKey := `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEID2nK2pCcGSbtS+U9jc2SCYxHWOo1eA4IR97bdif4+rx
-----END PRIVATE KEY-----`

	client := rasedi.NewClient(privateKey, secretKey)
	fmt.Printf("INSTALLED_OK %T\n", client)

	// Create Payment
	createPayload := rasedi.CreatePaymentPayload{
		Amount:                 "10200",
		Title:                  "Test Payment",
		Description:            "This is a test payment",
		Gateways:               []rasedi.Gateway{rasedi.GatewayCreditCard},
		RedirectURL:            "https://time.ir",
		CollectFeeFromCustomer: true,
		CollectCustomerEmail:   true,
		CollectCustomerPhoneNumber: false,
		CallbackURL:            "https://google.com",
	}

	createRes, err := client.CreatePayment(createPayload)
	var referenceCode string

	if err != nil {
		log.Printf("PAYMENT_CREATION_ERROR: %v", err)
		// Use a hardcoded one if creation failed for testing
		referenceCode = "0b0a8bce-bf3c-4fc4-993e-6179d95e9ece" 
	} else {
		createResJSON, _ := json.MarshalIndent(createRes, "", "  ")
		fmt.Printf("PAYMENT_CREATION_RESPONSE: %s\n", string(createResJSON))
		referenceCode = createRes.Body.ReferenceCode
	}

	// Get Payment
	if referenceCode != "" {
		getRes, err := client.GetPaymentByReferenceCode(referenceCode)
		if err != nil {
			log.Printf("GET_PAYMENT_ERROR: %v", err)
		} else {
			getResJSON, _ := json.MarshalIndent(getRes, "", "  ")
			fmt.Printf("GET_PAYMENT_RESPONSE: %s\n", string(getResJSON))
		}
	}

	// Cancel Payment
	if referenceCode != "" {
		cancelRes, err := client.CancelPayment(referenceCode)
		if err != nil {
			log.Printf("CANCEL_PAYMENT_ERROR: %v", err)
		} else {
			cancelResJSON, _ := json.MarshalIndent(cancelRes, "", "  ")
			fmt.Printf("CANCEL_PAYMENT_RESPONSE: %s\n", string(cancelResJSON))
		}
	}

	// Get After Cancel
	if referenceCode != "" {
		getAfterCancel, err := client.GetPaymentByReferenceCode(referenceCode)
		if err != nil {
			log.Printf("GET_AFTER_CANCEL_ERROR: %v", err)
		} else {
			getAfterCancelJSON, _ := json.MarshalIndent(getAfterCancel, "", "  ")
			fmt.Printf("GET_AFTER_CANCEL: %s\n", string(getAfterCancelJSON))
		}
	}
}
