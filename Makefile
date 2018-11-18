all: clean protos deps build install

clean:
	rm -rf include bin target readme.txt

deps:
	dep ensure

protos:
	wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip 
	unzip protoc.zip
	go get -u github.com/golang/protobuf/protoc-gen-go
	bin/protoc -I=rpc --go_out=plugins=grpc:rpc rpc/tcping.proto 

build:
	mkdir target
	go build -o target/tcping cmd/tcping/tcping.go
	go build -o target/tcpingd cmd/tcpingd/tcpingd.go
	go build -o target/tpctl cmd/tpctl/tpctl.go
	cp target/tcping $(GOPATH)/bin/tcping

install:
	sudo setcap cap_net_raw+ep $(GOPATH)/bin/tcping

.PHONY: clean deps build install protos
