package fastly

import (
	"context"
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"

	"github.com/pkg/errors"
)

// Factory creates a new usable instance of this secrets engine.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, errors.Wrap(err, "failed to create factory")
	}
	return b, nil
}

// backend is the actual backend.
type backend struct {
	*framework.Backend

	clientMutex sync.RWMutex
}

// Backend creates a new backend.
func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: []*framework.Path{
			// fastly/info
			&framework.Path{
				Pattern:      "info",
				HelpSynopsis: "Display information about this plugin",
				HelpDescription: `

Displays information about the plugin, such as the plugin version and where to
get help.

`,
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInfo,
				},
			},
			// fastly/config
			&framework.Path{
				Pattern:      "config",
				HelpSynopsis: "Configure fastly secret engine.",
				HelpDescription: `

Configure fastly secret engine.

`,
				Fields: map[string]*framework.FieldSchema{
					"username": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Fastly username",
						Default:     "",
					},
					"password": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Fastly password.",
						Default:     "",
					},
					"sharedSecret": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Fastly sharedSecret.",
						Default:     "",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   b.pathConfigRead,
					logical.UpdateOperation: b.pathConfigWrite,
				},
			},
			// fastly/generate
			&framework.Path{
				Pattern:      "generate",
				HelpSynopsis: "Generate and return a Fastly token",
				HelpDescription: `

Generate and return a Fastly token

`,
				Fields: map[string]*framework.FieldSchema{
					"scope": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "The scope for the token.",
						Default:     "",
					},
					"ttl": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "The ttl for the token to be valid.",
						Default:     "",
					},
					"service_id": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "The id of the service that the token is created for.",
						Default:     "",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathGenerate,
				},
			},
		},
	}

	return &b
}

func (b *backend) Close() {
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()
}

const backendHelp = `
The fastly secrets engine generates fastly tokens.
`
