# KoolSeal

A wrapper CLI to ease kubeseal secrets management on Kubernetes

## Download
```sh
go install github.com/ja1code/koolseal@latest
```

## Commands

`create` - Creates a new sealed secret file from a `.json` dictionary file

`update` - Get secret values from k8s cluster and generates a sealed file with patched values and/or new ones. 

`extract` - Get secrets from k8s, decodes and stores in a json file

## Examples

### Create a sealed secret file

Create a `.json` dictionary with the desired values
```json
// values.json file, at ~
{
  "DB_PORT": "3306",
  "DB_HOST": "local.svc",
  "DB_USER": "admin"
}
```

Call koolseal's `create` to generate a sealed secrets file
```shell
$ koolseal create -c cert.crt -ns default -n main-api -f ~/values.json ./secrets/main-api.secrets.sealed.yaml
```

- `-c` specifies the certificate to use
- `-ns` the namespace to be used on the new secret
- `-n` the name to be used on the new secret
- `-f` specifies the `.json` dictionary file location

The last argument is the destination where the sealed secrets file will be created.

### Update/Extend secrets

Create a `.json` dictionary with the desired updates and additions values
```json
// values.json file, at ~
{
  "DB_PORT": "3306",
  "DB_HOST": "local.svc",
  "DB_USER": "admin"
}
```

Call koolseal's `update` to generate a sealed secrets file
```shell
$ koolseal update -ns default -n main-api -f values.json -c cert.crt -p ./secrets/main-api.secret.sealed.yaml
```

- `-c` specifies the certificate to use
- `-n` the name of the secrets to be updated
- `-ns` the namespace of the secrets to be updated
- `-f` specifies the `.json` dictionary file location
- `-p` when in a git repository, you can automatically commit abd push the updates

The last argument is the destination where the updated sealed secrets file will be created, ideally you should inform the secret's current sealed secret file location, that way, koolseal will overwrite the previous secrets file with the update values

> You can also add/patch single values using the `--key` and `--value` flags and omitting the `--file`

### Extract

```shell
$ koolseal e -ns default -n main-api -f main-api.json
```

- `-ns` the namespace to be extracted
- `-n` the name and name to be extracted
- `-f` the file that will be created with the current secrets

This will generate a `.json` file in the following format:
```json
// values.json file, at ~
{
  "DB_PORT": "3306",
  "DB_HOST": "local.svc",
  "DB_USER": "admin"
}
```