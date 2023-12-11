# KoolSeal

A wrapper CLI to ease kubeseal secrets management on Kubernetes

## Download
```sh
$ go get github.com/ja1code/koolseal@latest
```

## Commands

`create` - Creates a new sealed secret file from a `.json` dictionary file

`update` - Get secret values from k8s cluster and generates a sealed file with patched values and/or new ones. 

`extract` - Get secrets from k8s, decodes and stores in a json file