# Sample-Credhub-KMS-Plugin
This is just a sample plugin. It contains the basic encrypt and decrypt functions using base64.
## Usage
```bash
mkdir -p $GOPATH/src/github.com/pivotal
git clone https://github.com/pivotal/sample-credhub-kms-plugin $GOPATH/src/github.com/pivotal/sample-credhub-kms-plugin
cd $GOPATH/src/github.com/pivotal/sample-credhub-kms-plugin
go build
./sample-credhub-kms-plugin /path/to/unix/socket
```

