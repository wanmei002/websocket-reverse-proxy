.PHONY: build-linux upload build-linux-client upload-client  buf

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy-linux main.go

upload:
	scp proxy-linux ubuntu@192.168.31.201:~/

build-linux-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy-client-linux client.go

upload-client:
	scp proxy-client-linux ubuntu@192.168.31.201:~/

buf:
	buf generate