package fastly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// FastlyToken is the response of Fastly create token API request
type FastlyToken struct {
	ID          string    `json:"id"`
	AccessToken string    `json:"access_token"`
	Name        string    `json:"name"`
	UserID      string    `json:"user_id"`
	ServiceID   string    `json:"service_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Scope       string    `json:"scope"`
	Services    []string  `json:"services"`
}

func (b *backend) pathGenerate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	config, err := b.config(ctx, req.Storage)

	formData := map[string][]string{
		"username": []string{config.Username},
		"password": []string{config.Password},
		"scope":    []string{"purge_all"},
		"name":     []string{"plugin_test"},

		// The Go url package doesn't add [] automatically when encoding arrays
		// This works, but looks hacky
		"services[]": []string{"Xj6ZldKCnW2gmTix97F1U", "23MQEr22Ux4U7rW4IRZUXz"}, // Sandbox services
	}

	token, err := CreateFastlyToken(config.SharedSecret, formData)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token": token,
		},
	}, nil
}

// CreateFastlyToken uses Fastly API to create an API token
func CreateFastlyToken(otp string, formData map[string][]string) (*FastlyToken, error) {
	// Make sure the token expires
	formData = ensureExpirationInParams(formData)

	// Prepare request
	form := url.Values(formData)
	req, err := http.NewRequest(
		"POST",
		"https://api.fastly.com/tokens",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	// Add OTP to header if provided
	if len(otp) > 0 {
		req.Header.Set("Fastly-OTP", otp)
	}

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read response body into a buffer
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	// Print response status and body if the status is unexpected
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("%d %s %s", resp.StatusCode, resp.Status, buf.String())
	}

	// Parse response JSON
	token := &FastlyToken{}
	if err := json.Unmarshal(buf.Bytes(), token); err != nil {
		return nil, err
	}

	return token, nil
}

func ensureExpirationInParams(formData map[string][]string) map[string][]string {
	// This adds expiration to create token POST form if one doesn't exist
	// Default TTL 5 minutes for now
	defaultExpiration := time.Now().Add(5 * time.Minute).UTC().Format("2006-01-02T15:04:05+00:00")

	field := formData["expires_at"]
	if field == nil || len(field) == 0 {
		formData["expires_at"] = []string{defaultExpiration}
	}

	return formData
}
