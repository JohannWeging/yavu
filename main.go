package main

import (
	"time"

	"github.com/JohannWeging/logerr"
	vault "github.com/hashicorp/vault/api"
	"github.com/juju/errors"
	"github.com/luzifer/rconfig"
	log "github.com/sirupsen/logrus"
)

var config = struct {
	CheckInterval   time.Duration `flag:"check-interval" default:"60s" description:"interval to check the vault"`
	SourceVaultAddr string        `flag:"source-vault-addr" description:"vault address to read the unseal key from"`
	SourceRoleID    string        `flag:"source-role-id" description:"role id to authenticate at the source vault"`
	SourceSecretID  string        `flag:"source-secret-id" description:"secret id to authenticate at the source vault"`
	TargetVaultAddr string        `flag:"target-vault-addr" description:"vault address of the instance to unseal"`
	UnsealTokenPath string        `flag:"unseal-token-path" description:"where to read the unseal tokens from"`
	UnsealTokenKeys []string      `flag:"unseal-token-keys" description:"list unseal keys in the secret path"`
}{}

func main() {
	if err := run(); err != nil {
		fields := logerr.GetFields(err)
		log.WithFields(fields).WithError(err).Fatal("failed to start yavu")
	}
}

func run() error {
	rconfig.AutoEnv(true)
	if err := rconfig.Parse(&config); err != nil {
		return errors.Annotate(err, "failed to parse config")
	}
	return check()
}

func readUnsealTokens() (tokens []string, err error) {
	logFields := logerr.Fields{"vault_addr": config.SourceVaultAddr}
	logerr.DeferWithFields(&err, logFields)
	errors.DeferredAnnotatef(&err, "failed to read unseal tokens")
	c := &vault.Config{
		Address: config.SourceVaultAddr,
	}
	client, err := vault.NewClient(c)
	if err != nil {
		return nil, errors.Annotate(err, "failed to create vault client")
	}
	if config.SourceRoleID == "" {
		return nil, errors.New("no source sole id provided")
	}

	data := map[string]interface{}{"role_id": config.SourceRoleID}
	if config.SourceSecretID != "" {
		data["secret_id"] = config.SourceSecretID
	}

	loginSecret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return nil, errors.Annotate(err, "failed to fetch authentication token")
	}
	if loginSecret.Auth == nil {
		return nil, errors.New("authentication token is nil")
	}
	client.SetToken(loginSecret.Auth.ClientToken)
	logFields["unseal_token_path"] = config.UnsealTokenPath
	secret, err := client.Logical().Read(config.UnsealTokenPath)
	if err != nil {
		return nil, errors.Annotate(err, "failed to read unseal tokens")
	}
	if secret == nil {
		return nil, errors.New("secret is nil")
	}

	if secret.Data == nil {
		return nil, errors.Errorf("secret data is nil")
	}
	for _, key := range config.UnsealTokenKeys {
		logFields["unseal_key"] = key
		data, ok := secret.Data[key]
		if !ok {
			err = errors.New("secret is missing token key")
			return nil, logerr.WithField(err, "key", key)
		}
		token, ok := data.(string)
		if !ok {
			err = errors.New("secret data is not of type string")
			return nil, logerr.WithField(err, "key", key)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func check() (err error) {
	c := &vault.Config{Address: config.TargetVaultAddr}
	client, err := vault.NewClient(c)
	if err != nil {
		err = logerr.WithField(err, "vault_addr", config.TargetVaultAddr)
		return errors.Annotate(err, "failed to create target vault client")
	}

	for _ = range time.Tick(config.CheckInterval) {
		err = unseal(client)
		if err != nil {
			fields := logerr.GetFields(err)
			fields["vault_addr"] = config.TargetVaultAddr
			log.WithFields(fields).WithError(err).Error("failed to unseal the vault")
		}
	}
	return nil
}

func unseal(client *vault.Client) (err error) {
	errors.DeferredAnnotatef(&err, "failed to unseal instance")

	status, err := client.Sys().SealStatus()
	if err != nil {
		return errors.Annotate(err, "failed to get seal status")
	}
	if !status.Initialized {
		return errors.New("vault not initialized")
	}
	if !status.Sealed {
		return nil
	}
	tokens, err := readUnsealTokens()
	if err != nil {
		return err
	}
	for _, t := range tokens {
		resp, err := client.Sys().Unseal(t)
		if err != nil {
			return errors.Annotate(err, "failed to send unseal token")
		}
		if !resp.Sealed {
			break
		}
	}
	return nil
}
