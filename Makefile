build:
	go build -o bin/cdp cmd/cdp/main.go

build-proto:
	buf mod update apis
	buf lint
	buf build
	buf generate


test-cmd:
	grpcurl  -plaintext 127.0.0.1:50052 dataplugin.v1alpha.DataPluginService.Healthiness
	grpcurl  -plaintext 127.0.0.1:50052 dataplugin.v1alpha.DataPluginService.Registration
