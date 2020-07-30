# secrets provider

## Init
- `go mod download`

## Build
- `make build`

## Run:
- `./bin/secrets-provider_<build_date> -port 8090 -persist`
(persist flag creates secrets.db in current dir)
- enter crypt key (e.g. "qwerty")

## Usage
- `curl -X POST -H "X-Crypt-Key: qwerty" -d '{"secret": "1234"}' localhost:8090/secret1` - create or update
- `curl -X GET -H "X-Crypt-Key: qwerty" localhost:8090/secret1` - get
- `curl -X DELETE -H "X-Crypt-Key: qwerty" localhost:8090/secret1` - delete
