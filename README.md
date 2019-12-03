# vault-fastly-secret-engine

This plugin will allow you to create a secret backend that will use the Fastly API to generate dynamic short lived Fastly token.  Usage can be restricted using the highly customizable Vault ACL system.

### Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Register the plugin with the catalog

    ```text
    $ SHASUM=$(shasum -a 256 vault-fastly-secret-engine | cut -d " " -f1)
    $ vault write sys/plugins/catalog/vault-fastly-secret-engine sha_256="$SHASUM" command="vault-fastly-secret-engine" 
    Success! Data written to: sys/plugins/catalog/vault-fastly-secret-engine
    ```

1. Enable the Fastly secrets engine:

    ```text
    $ vault secrets enable -path="fastly" -plugin-name="vault-fastly-secret-engine" plugin
    Success! Enabled the vault-fastly-secret-engine plugin at: fastly/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

1. Configure the backend with user credentials that will be able to interact with the Fastly API and create tokens.

    ```text
    $ vault write fastly/config username="sam" password="test" sharedSecret="123"
    Success! Data written to: fastly/config
    ```

    The `sharedSecret` corresponds to the shared secret key produced by fastly when configuring MFA login.  This will be used to generate the Fastly tokens.

### Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate tokens.

1. Generate a new fastly token by writing to the  `/fastly/generate` endpoint with the
scope of the desired token as well as the service ID:

    ```text
    $ vault write fastly/generate scope="global" service_id="Xj62345gmTix9gh67U"
    Key      Value
    ---      -----
    token    d118a65cdfe314202cf969e1fb2e8afc
    ```

    *NOTE* you can provide multiple servie IDs by using a comma delimited string.

    ```text
    $ vault write fastly/generate scope="global" service_id="Xj62345gmTix9gh67U,45MDE6457BT4IRZdf7z"
    Key      Value
    ---      -----
    token    f2732f475773ab0d0bce1cd371d72b48
    ```

    Using ACLs, it is possible to restrict the type of tokens that can be generated.  Any combination of scope and service ID can be used

## Local Development

### Build the code

```bash
GOOS=linux GOARCH=amd64 go build
docker build -t vault-plugin .
docker run --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:1234' -p 1234:1234 vault-plugin
```

### Configure the local vault

In a second terminal window...

```bash
export VAULT_ADDR='http://0.0.0.0:1234'
vault login myroot
SHASUM=$(shasum -a 256 vault-fastly-secret-engine | cut -d " " -f1)
vault write sys/plugins/catalog/vault-fastly-secret-engine   sha_256="$SHASUM"   command="vault-fastly-secret-engine"
vault secrets enable -path="fastly" -plugin-name="vault-fastly-secret-engine" plugin
vault write fastly/config username="sam" password="test" sharedSecret="123"
```
