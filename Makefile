TARGET_FILE:=mock-http
BUILD_DIR=.build
COVER_PROFILE_FILE="${BUILD_DIR}/go-cover.tmp"
GOBIN := $(GOPATH)/bin

clean:
	rm -rf $(TARGET_FILE) $(BUILD_DIR)

clean-test:
	@go fmt ./...
	@go clean -testcache

test: clean-test
	go test -p 1 ./...

cover-html: mk-build-dir clean-test
	go test -p 1 -coverprofile=${COVER_PROFILE_FILE} ./... ; echo
	go tool cover -html=${COVER_PROFILE_FILE}
	$(MAKE) badge

build-deps:
	@go mod tidy

build: build-deps test
	CGO_ENABLED=1 go build -tags "json1 fts5 foreign_keys math_functions" -o $(TARGET_FILE)

install: build
	cp $(TARGET_FILE) $(GOBIN)/$(TARGET_FILE)

mk-build-dir: build
	mkdir -p .build
