# yavu
Yet Another Vault Unsealer, unseals independent vault instances.

But why?<br/>
yavu reads the the unseal key from vault *A* to unseal vault *B*.<br\>
Vault *A* must be unsealed to unseal vault *B*.<br/>
The unseal key is "not stored" in memory, every unseal try of vault B reads the unseal token from vault A.<br/>
The vault infrasturcture can be sealed without the requriment to kill the unseal processes.<br/>

## Usage
```
yavu --source-vault-addr=http://vault.A 
  --source-role-id=foo-bar  
  --unseal-token-path=secret/vault-b-unseal
  --unseal-token-keys=key1,key2
  --target-vault-addr=http://vault.B
```
Or with Docker
```
docker run 
  -e SOURCE_VAULT_ADDR=http://vault.A
  -e SOURCE_ROLE_ID=foo-bar
  -e UNSEAL_TOKEN_PATH=secret/vault-b-unseal
  -e UNSEAL_TOKEN_KEYS=key1,key2
  -e TARGET_VAULT_ADDR=http://vault.B
  johannweging/yavu
```
Run reference 
```
Usage of yavu:
      --check-interval duration     interval to check the vault (default 1m0s)
      --source-role-id string       role id to authenticate at the source vault
      --source-secret-id string     secret id to authenticate at the source vault
      --source-vault-addr string    vault address to read the unseal key from
      --target-vault-addr string    vault address of the instance to unseal
      --unseal-token-keys strings   list unseal keys in the secret path
      --unseal-token-path string    where to read the unseal tokens from
```
