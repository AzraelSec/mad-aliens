OUTPUT_DIR=./bin
CMD_DIR=./cmd

build: build_cli

build_cli:
	GOARCH=amd64 GOOS=darwin go build -o ${OUTPUT_DIR}/cli/cli-darwin ${CMD_DIR}/cli/main.go
	GOARCH=amd64 GOOS=linux go build -o ${OUTPUT_DIR}/cli/cli-linux ${CMD_DIR}/cli/main.go
	GOARCH=amd64 GOOS=windows go build -o ${OUTPUT_DIR}/cli/cli-windows ${CMD_DIR}/cli/main.go

test:
	go test ./...

clean:
	go clean
	rm -r ${OUTPUT_DIR}

