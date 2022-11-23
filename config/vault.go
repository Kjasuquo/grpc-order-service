package config

import (
	"context"
	"encoding/json"

	// "fmt"
	vault "github.com/hashicorp/vault/api"
	// "github.com/hashicorp/vault/api/auth/approle"
	"log"
)

type VaultParameters struct {
	// connection parameters
	address             string
	approleRoleID       string
	approleSecretIDFile string

	apiKeyPath              string
	apiKeyMountPath         string
	apiKeyField             string
	databaseCredentialsPath string
}

type Vault struct {
	client     *vault.Client
	parameters VaultParameters
}

func VaultSecrets(vaultAdd, vaultToken, secretPath string) (*Config, error) {
	vaultConfig := vault.DefaultConfig()
	// vaultConfig.Address = viper.GetString("VAULT_ADDR")
	vaultConfig.Address = vaultAdd
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	client.SetToken(vaultToken)
	// client.SetToken(viper.GetString("VAULT_AUTH_TOKEN"))

	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("secret").Get(context.Background(), secretPath)

	// fmt.Printf("from vault:: ,%v \n", secret.Data)
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	j, err := json.Marshal(secret.Data)

	if err != nil {
		log.Fatalf("unable to marshal secrets: %v", err)
	}

	// fmt.Printf("j data %T, %v /n", j, j)

	config := &Config{}

	err = json.Unmarshal(j, config)
	if err != nil {
		log.Fatalf("unable to parse secrets: %v", err)
	}
	// fmt.Printf("j data %T, %v /n", j, j)

	return config, nil
}
