.PHONY: build
build:
	go build

.PHONY: test
test: execute-test
	diff test/expected.txt test/actual/service_code_generated.txt && echo "PATH"

.PHONY: execute-test
execute-test: build clean
	mkdir test/actual && \
	protoc --plugin=./protoc-gen-service --service_out=configPath=./test/excludes.yml:test/actual --proto_path=./test/proto -I ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.7 ./test/proto/service.proto && \
	mv test/actual/service_code_generated.go test/actual/service_code_generated.txt

.PHONY: clean
clean:
	if [ -e test/actual ]; then rm -rf test/actual; fi
