FROM golang:1.17

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git bash ca-certificates build-essential curl unzip

RUN curl -OL https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
    unzip protoc-3.6.1-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm protoc-3.6.1-linux-x86_64.zip && \
    rm -rf protoc3

RUN go get github.com/golang/protobuf/protoc-gen-go && \
    go get github.com/jucardi/protoc-go-inject-tag && \
    go get golang.org/x/tools/cmd/goimports && \
    go get github.com/jucardi/goimports-blank-rm
