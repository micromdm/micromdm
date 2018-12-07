# MicroMDM - Abacus
In addition, to official README of the source project: https://github.com/micromdm/micromdm

## Requirements

Make sure, go version 1.11 or newer is installed:
`go version`

For testing, install `goose`
`brew install goose`

## Compiling
Run build in root directory of the project
`make build`

Will generate our executables in 
```
./build/darwin
./build/linux
```

Will also generate a `docker-compose-dev.yaml` file, for further configuration.

Choose your system (darwin/linux) and go to the directory.

```
# Executable for configuration and setup
./build/(darwin|linux)/mdmctl

# Executable for running the server
./build/(darwin|linux)/micromdm
```

## Testing
### Mysql
#### Setup test account:
* username: micromdm
* password: micromdm
```
mysql -u root
mysql> CREATE USER 'micromdm'@'localhost' IDENTIFIED BY 'micromdm';
mysql> GRANT ALL PRIVILEGES ON *.* TO 'micromdm'@'localhost' IDENTIFIED BY 'micromdm';
```
#### Run tests
```
make db-mysql-migrate-test
make db-mysql-test
```

## Setup
### mdmctl 
#### Configure later MDM Service
```
./mdmctl config set \
    -api-token secret \
    -name mdmexample \
    -server-url https://mdm.abacus.ch/
./mdmctl config switch -name mdmexample
```

#### Assign Apple Push Certificate
To assign an Apple Push Certificate, start the server first (no Mysql database connection required, we won't store the certificate in the Mysql database, but locally in a document store.)
```
sudo ./micromdm serve \
    -api-key secret \
    -tls-cert ./fullchain.pem \
    -tls-key ./privkey.pem \
    -server-url https://mdm.abacus.ch/
```

Now, when the server is running, add the Push certificate, from a second Console to the server.
```
./mdmctl mdmcert upload \
    -password secret \
    -cert ./mdm-certificates/MDM_\ Abacus\ Research\ AG_Certificate.pem \
    -private-key ./mdm-certificates/PushCertificatePrivateKey.key
```

### micromdm
After configuring the MDM Service, run it.

```
sudo ./micromdm serve \
    -api-key secret \
    -tls-cert ./fullchain.pem \
    -tls-key ./privkey.pem \
    -server-url https://mdm.abacus.ch/ \
    -command-webhook-url http://127.0.0.1:5000/webhook \
    -mysql-username micromdm \
    -mysql-password micromdm \
    -mysql-database micromdm_test \
    -mysql-host 127.0.0.1 \
    -mysql-port 3306
```
