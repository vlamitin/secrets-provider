BUILD_DATE=`date +%d_%m_%YT%H_%M_%S`
CGO_ENABLED=1
WINDOWS_ENV=env GOOS=windows GOARCH=amd64
LINUX_ENV=env GOOS=linux GOARCH=amd64

run:
	go run cmd/app/main.go -port 8090

run-persist:
	go run cmd/app/main.go -port 8090 -persist

# Builds okko_films binaries for amd64 linux and amd64 windows
build:
	${WINDOWS_ENV} go build -o bin/secrets-provider_${BUILD_DATE}.exe cmd/app/main.go
	${LINUX_ENV} go build -o bin/secrets-provider_${BUILD_DATE} cmd/app/main.go

# Cleans our project: deletes binaries
clean:
	rm -f bin/*
