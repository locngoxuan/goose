.PHONY: dist

BUILD_NAME=goose

clean:
	rm -rf ./bin

vendor:
	go mod vendor && go mod tidy

dev:
	@mkdir -p ./bin
	@rm -f ./bin/*
	env CGO_ENABLED=1 go build -tags='no_mssql no_mysql no_sqlite3' -o ./bin/goose ./cmd/goose

dist:
	@mkdir -p ./bin
	@rm -f ./bin/*
	env CGO_ENABLED=0 go build -tags='no_mssql no_mysql no_sqlite3' -ldflags="-s -w" -o ./bin/goose ./cmd/goose
	env CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -tags='no_mssql no_mysql no_sqlite3' -ldflags="-s -w" -o ./bin/goose-darwin       ./cmd/goose
	env CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags='no_mssql no_mysql no_sqlite3' -ldflags="-s -w" -o ./bin/goose-linux        ./cmd/goose
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags='no_mssql no_mysql no_sqlite3' -o ./bin/goose-wins.exe  ./cmd/goose

compress:
	upx --brute ./bin/${BUILD_NAME}-linux
	upx --brute ./bin/${BUILD_NAME}-darwin