package fastly

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type fastlyConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SharedSecret string `json:"sharedSecret"`
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	config, err := b.config(ctx, req.Storage)

	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &fastlyConfig{}
	}

	if err := config.Update(data); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("could not update config: %v", err)), nil
	}

	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Invalidate existing clients so they read the new configuration
	b.Close()

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	resp := make(map[string]interface{})

	if v := config.Username; v != "" {
		resp["username"] = v
	}
	if v := config.Password; v != "" {
		resp["password"] = v
	}
	if v := config.SharedSecret; v != "" {
		resp["sharedSecret"] = v
	}

	return &logical.Response{
		Data: resp,
	}, nil
}

func (config *fastlyConfig) Update(data *framework.FieldData) error {
	username := data.Get("username").(string)
	if len(username) > 0 {
		config.Username = username
	}

	password := data.Get("password").(string)
	if len(password) > 0 {
		config.Password = password
	}

	sharedSecret := data.Get("sharedSecret").(string)
	if len(sharedSecret) > 0 {
		config.SharedSecret = strings.ToUpper(sharedSecret)
	}

	return nil
}

func (b *backend) config(ctx context.Context, s logical.Storage) (*fastlyConfig, error) {
	config := &fastlyConfig{}
	entry, err := s.Get(ctx, "config")

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}
