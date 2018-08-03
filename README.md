# vault-fastly-secret-engine
Vault secret engine for Fastly

## To run locally

```bash
GOOS=linux GOARCH=amd64 go build
docker build -t vault-plugin .
docker run --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:1234' -p 1234:1234 vault-plugin
```

In a second terminal window...

```bash
export VAULT_ADDR='http://0.0.0.0:1234'
vault login myroot
SHASUM=$(shasum -a 256 vault-fastly-secret-engine | cut -d " " -f1)
vault write sys/plugins/catalog/vault-fastly-secret-engine   sha_256="$SHASUM"   command="vault-fastly-secret-engine"
vault secrets enable -path="fastly" -plugin-name="vault-fastly-secret-engine" plugin
```

At this point the sercret engine is enabled and you can interact with it.  To configure the engine: 

```bash
vault write fastly/config username="sam" password="test" sharedSecret="123"
```

You can view the config with: 

```bash
vault read fastly/config
```

You can generate a token with:
```bash
vault write fastly/generate scope="" service_id=""
```
