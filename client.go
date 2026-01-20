package rasedi

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const (
	APIBaseURL      = "https://api.rasedi.com"
	UpstreamVersion = 1
)

type Client struct {
	authenticator *Auth
	httpClient    *http.Client
	isTest        bool
	publicKeys    []PublicKey
}

// NewClient initializes a new Rasedi SDK client.
// It determines the environment (test/live) based on the presence of "test" in the secretKey.
func NewClient(privateKey, secretKey string) *Client {
	client := &Client{
		authenticator: NewAuth(privateKey, secretKey), // secretKey is used as keyID/passphrase based on python logic
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}
	client.isTest = client.checkIsTest(secretKey)
	return client
}

func (c *Client) checkIsTest(secret string) bool {
	return strings.Contains(secret, "test")
}

func (c *Client) call(path, method string, requestBody interface{}) (map[string]interface{}, error) {
	env := "live"
	if c.isTest {
		env = "test"
	}
	relativeURL := fmt.Sprintf("/v%d/payment/rest/%s%s", UpstreamVersion, env, path)
	fullURL := fmt.Sprintf("%s%s", APIBaseURL, relativeURL)

	signature, err := c.authenticator.MakeSignature(method, relativeURL)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-signature", signature)
	req.Header.Set("x-id", c.authenticator.GetKeyID())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}


	// Try to unmarshal as JSON, if fails, keep as string/nil inside the map if needed,
	// but the return type expects map[string]interface{}.
	// The Python SDK returns "body" as parsed JSON or text.
	// Here we will try to unmarshal into a generic map or list.
	// But our specific methods expect specific structs.
	// Let's return the simplified response structure:
	// map[string]interface{} with keys "body", "headers", "statusCode"

	headerMap := make(map[string]string)
	for k, v := range resp.Header {
		headerMap[k] = strings.Join(v, ",")
	}

	result := map[string]interface{}{
		"headers":    headerMap,
		"statusCode": resp.StatusCode,
	}

	if json.Valid(bodyBytes) {
		var decodedBody interface{}
		_ = json.Unmarshal(bodyBytes, &decodedBody)
		result["body"] = decodedBody
	} else {
		result["body"] = string(bodyBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 209 {
		return result, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return result, nil
}

// Public Methods

// GetPublicKeys fetches the list of public keys used for signature verification.
func (c *Client) GetPublicKeys() (*PublicKeysResponse, error) {
	resp, err := c.call("/get-public-keys", "GET", nil)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := json.Marshal(resp["body"])
	var publicKeys []PublicKey
	if err := json.Unmarshal(bodyBytes, &publicKeys); err != nil {
		return nil, err
	}

	c.publicKeys = publicKeys

	return &PublicKeysResponse{
		Body:       publicKeys,
		Headers:    resp["headers"].(map[string]string),
		StatusCode: resp["statusCode"].(int),
	}, nil
}

// CreatePayment initiates a new payment request.
func (c *Client) CreatePayment(payload CreatePaymentPayload) (*CreatePaymentResponse, error) {
	resp, err := c.call("/create", "POST", payload)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := json.Marshal(resp["body"])
	var responseBody CreatePaymentResponseBody
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		return nil, err
	}

	return &CreatePaymentResponse{
		Body:       responseBody,
		Headers:    resp["headers"].(map[string]string),
		StatusCode: resp["statusCode"].(int),
	}, nil
}

// GetPaymentByReferenceCode retrieves the status and details of a payment using its reference code.
func (c *Client) GetPaymentByReferenceCode(referenceCode string) (*PaymentDetailsResponse, error) {
	resp, err := c.call(fmt.Sprintf("/status/%s", referenceCode), "GET", nil)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := json.Marshal(resp["body"])
	var responseBody PaymentDetailsResponseBody
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		return nil, err
	}

	return &PaymentDetailsResponse{
		Body:       responseBody,
		Headers:    resp["headers"].(map[string]string),
		StatusCode: resp["statusCode"].(int),
	}, nil
}

// CancelPayment cancels a pending payment using its reference code.
func (c *Client) CancelPayment(referenceCode string) (*CancelPaymentResponse, error) {
	resp, err := c.call(fmt.Sprintf("/cancel/%s", referenceCode), "PATCH", nil)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := json.Marshal(resp["body"])
	var responseBody CancelPaymentResponseBody
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		return nil, err
	}

	return &CancelPaymentResponse{
		Body:       responseBody,
		Headers:    resp["headers"].(map[string]string),
		StatusCode: resp["statusCode"].(int),
	}, nil
}

// Verify validates the webhook signature and returns the decoded payload.
// It automatically fetches public keys if they are not already cached.
func (c *Client) Verify(payload VerifyPayload) (map[string]interface{}, error) {
	if len(c.publicKeys) == 0 {
		if _, err := c.GetPublicKeys(); err != nil {
			return nil, fmt.Errorf("failed to fetch public keys: %v", err)
		}
	}

	targetKeyID := payload.KeyID
	var targetKey *PublicKey
	for _, k := range c.publicKeys {
		if k.ID == targetKeyID {
			targetKey = &k
			break
		}
	}

	if targetKey == nil {
		// Retry fetching public keys
		if _, err := c.GetPublicKeys(); err != nil {
			return nil, fmt.Errorf("failed to fetch public keys on retry: %v", err)
		}
		for _, k := range c.publicKeys {
			if k.ID == targetKeyID {
				targetKey = &k
				break
			}
		}
		if targetKey == nil {
			return nil, fmt.Errorf("internal server error: public key not found")
		}
	}

	if payload.Content == nil {
		return nil, fmt.Errorf("internal server error: empty content")
	}
	content := *payload.Content

	parts := strings.Split(content, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerB64, payloadB64, sigB64 := parts[0], parts[1], parts[2]
	signedData := []byte(fmt.Sprintf("%s.%s", headerB64, payloadB64))

	sigRaw, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		// Try standard decoding if raw fails (though JWS usually uses URL encoding)
		sigRaw, err = base64.URLEncoding.DecodeString(sigB64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode signature: %v", err)
		}
	}

	block, _ := pem.Decode([]byte(targetKey.Key))
	if block == nil {
		return nil, fmt.Errorf("failed to parse public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPub, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unsupported public key type for ES signature verification")
	}

	keySizeBytes := (ecdsaPub.Curve.Params().BitSize + 7) / 8
	if len(sigRaw) != 2*keySizeBytes {
		// This check is strict, some implementations might have different signature formats
		// But valid R/S split usually follows this for P-curves.
		// Let's assume standard P-256/384/521 concatenation.
		// However, Python code manually constructs ASN.1. Go's VerifyASN1 needs ASN.1.
		// We have raw R||S, so we need to convert to ASN.1 or use Verify with R and S big.Ints.
		// Go's ecdsa.Verify expects R and S.
	}

	r := new(big.Int).SetBytes(sigRaw[:keySizeBytes])
	s := new(big.Int).SetBytes(sigRaw[keySizeBytes:])

	hasher := sha512.New()
	hasher.Write(signedData)
	hash := hasher.Sum(nil)

	valid := ecdsa.Verify(ecdsaPub, hash, r, s)

	if !valid {
		return nil, fmt.Errorf("signature verification failed")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payloadJSON interface{}
	if err := json.Unmarshal(payloadBytes, &payloadJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	return map[string]interface{}{
		"body":        payloadJSON,
		"headers":     map[string]string{},
		"status_code": 200,
	}, nil
}
