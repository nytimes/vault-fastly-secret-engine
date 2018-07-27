package fastly

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"github.com/nytm/vault-fastly-secret-engine/version"
)

func (b *backend) pathInfo(_ context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"commit":  version.GitCommit,
			"version": version.Version,
		},
	}, nil
}
