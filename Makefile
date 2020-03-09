.PHONY: dist

BUILD_NAME=goose

clean:
	rm -rf ./bin

dist:
	@mkdir -p ./bin
	@rm -f ./bin/*
	go mod vendor && go mod tidy
	env CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/goose-darwin       ./cmd/goose
	env CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/goose-linux        ./cmd/goose
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/goose-wins.exe  ./cmd/goose

compress:
	upx --brute ./bin/${BUILD_NAME}-linux
	upx --brute ./bin/${BUILD_NAME}-darwin


.PHONY: vendor
vendor:
	mv _go.mod go.mod
	mv _go.sum go.sum
	GO111MODULE=on go build -o ./bin/goose ./cmd/goose
	GO111MODULE=on go mod vendor && GO111MODULE=on go mod tidy
	mv go.mod _go.mod
	mv go.sum _go.sum
