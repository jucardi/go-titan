all: check deps test build

check: format vet

format:
	@echo "formatting files..."
	@echo "  > getting goimports" && go install golang.org/x/tools/cmd/goimports@latest
	@echo "  > getting goimports-blank-rm" && go install github.com/jucardi/goimports-blank-rm@latest
	@echo "  > executing goimports-blank-rm" && goimports-blank-rm . 1>/dev/null 2>/dev/null
	@echo "  > executing goimports" && goimports -w -l $(shell find . -type f -name '*.go' -not -path "./vendor/*") 1>/dev/null 2>/dev/null || true
	@echo "  > executing gofmt" && gofmt -s -w -l $(shell find . -type f -name '*.go' -not -path "./vendor/*") 1>/dev/null

vet:
	@echo "vetting..."
	@go vet -mod=vendor ./...

templates: protoc format

protoc:
	@echo "generating protobuf..."
	@go get google.golang.org/protobuf/protoc-gen-go
	@GO111MODULE=off go get github.com/jucardi/protoc-go-inject-tag
	@GO111MODULE=off go install github.com/jucardi/protoc-go-inject-tag
	@protoc -I=$(PWD)/net/errorx --go_out=$(PWD)/net/errorx $(PWD)/net/errorx/error.proto
	@protoc-go-inject-tag --input "$(PWD)/net/errorx/*.pb.go" --cleanup -x yaml -x gorm -x bson

deps: protoc
	@echo "installing dependencies..."
	@go get ./...
	@go mod tidy
	@go mod vendor

test:
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@go test -mod=vendor ./... -v -coverprofile test-artifacts/cover.out
	@go tool cover -func test-artifacts/cover.out

build:
	@echo "building..."
	@go build -mod=vendor ./...

docker-cleanup:
	@docker-compose -f build-compose.yml down 1>/dev/null

docker-build: docker-cleanup
	@docker-compose -f build-compose.yml up --abort-on-container-exit --exit-code-from builder builder
	@docker-compose -f build-compose.yml down -v