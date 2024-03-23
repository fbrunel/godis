TARGET_CLI=godis-cli
TARGET_SRV=godis-server

build:
	go build -o bin/$(TARGET_CLI) cmd/cli/godis-cli.go
	go build -o bin/$(TARGET_SRV) cmd/srv/godis-server.go

clean:
	go clean
	rm -f bin/$(TARGET_CLI)
	rm -f bin/$(TARGET_SRV)
