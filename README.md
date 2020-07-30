# secrets provider

## Run:
- `go run cmd/app/main.go -port 8090`
- enter crypt key (e.g. "qwerty")

## Usage
- `curl -X POST -H "X-Crypt-Key: qwerty" -d '{"secret": "1234"}' localhost:8090/secret1`
- `curl -X GET -H "X-Crypt-Key: qwerty" localhost:8090/secret1`
- `curl -X DELETE -H "X-Crypt-Key: qwerty" localhost:8090/secret1`
