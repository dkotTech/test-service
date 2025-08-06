.PHONY: all
all:
	protoc -I/usr/local/include -I. -I${GOPATH}/src --go_out=./ ./balances/grpc/*.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src --go_out=./ --go-grpc_out=./ ./balances/grpc/*.proto
