#!make
install: protogen
	go get ./...

protogen:
	protoc -I/usr/local/include -I. \
	  --go_out=plugins=micro:$(GOPATH)/src/github.com/reversTeam/fizzbuzz-golang/src/client \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:. \
	src/**/protobuf/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	src/**/protobuf/*.proto
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	src/**/protobuf/*.proto

clean:
	rm src/**/protobuf/*.pb.go || true
	rm src/**/protobuf/*.pb.gw.go || true
	rm src/**/protobuf/*.swagger.json || true

run:
	go run gateway.go

build:
	GOOS=linux GOARCH=amd64 go build -o gateway ./main.go
	GOOS=linux GOARCH=amd64 go build -o client ./src/client/main.go
	docker build -t triviere42/fizzbuzz-golang .
	docker push triviere42/fizzbuzz-golang