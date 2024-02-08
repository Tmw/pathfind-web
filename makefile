wasm:
	@pushd go && GOARCH=wasm GOOS=js go build -o ../bridge.wasm .
	
