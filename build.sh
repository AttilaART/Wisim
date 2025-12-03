#! /bin/env bash

cd wisim-nouveau &&
	npm i &&
	cd .. &&
	cd wisim-backend/server &&
	$(
		GOOS=windows GOARCH=amd64 go build -o wisimserver-windows-amd64.exe wisimserver.go &
		GOOS=darwin GOARCH=amd64 go build -o wisimserver-darwin-amd64 wisimserver.go &
		GOOS=darwin GOARCH=arm64 go build -o wisimserver-darwin-arm64 wisimserver.go &
		GOOS=linux GOARCH=amd64 go build -o wisimserver-linux-amd64 wisimserver.go
	)
cd ../../ &&
	cd wisim-backend/wasm &&
	GOOS=js GOARCH=wasm go build -o main.wasm main.go &&
	cd ../../ &&
	mv wisim-backend/wasm/main.wasm wisim-nouveau/static/main.wasm &&
	mv wisim-backend/server/wisimserver-windows-amd64.exe wisim-nouveau/static &&
	mv wisim-backend/server/wisimserver-darwin-amd64 wisim-nouveau/static &&
	mv wisim-backend/server/wisimserver-darwin-arm64 wisim-nouveau/static &&
	mv wisim-backend/server/wisimserver-linux-amd64 wisim-nouveau/static &&
	cd wisim-nouveau &&
	npm run build &&
	cd .. &&
	echo "building binary" &&
	$(
		GOOS=windows GOARCH=amd64 go build -o build/wisim-windows-amd64.exe main.go &
		GOOS=darwin GOARCH=amd64 go build -o build/wisim-darwin-amd64 main.go &
		GOOS=darwin GOARCH=arm64 go build -o build/wisim-darwin-arm64 main.go &
		GOOS=linux GOARCH=amd64 go build -o build/wisim-linux-amd64 main.go
	) &&
	echo "done"
