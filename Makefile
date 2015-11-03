configure:
	gb vendor update --all

build:
	gofmt -w src/wssrv
	go tool vet src/wssrv/*.go
	gb test
	gb build