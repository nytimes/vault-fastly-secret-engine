package fastly

import (
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

// validateFields verifies that no bad arguments were given to the request.
func validateFields(req *logical.Request, data *framework.FieldData) error {
	var unknownFields []string
	for k := range req.Data {
		if _, ok := data.Schema[k]; !ok {
			unknownFields = append(unknownFields, k)
		}
	}

	if len(unknownFields) > 0 {
		// Sort since this is a human error
		sort.Strings(unknownFields)

		return fmt.Errorf("unknown fields: %q", unknownFields)
	}

	return nil
}

func generateTOTPCode(key string) (string, error) {

	// Generate password using totp library
	totpToken, err := totplib.GenerateCodeCustom(key, time.Now(), totplib.ValidateOpts{
		Period:    30,
		Digits:    6,
		Algorithm: otplib.AlgorithmSHA1,
	})

	if err != nil {
		return "", err
	}

	// Return the secret
	return totpToken, nil

}
