BINARY=gprobe
VERSION=$(shell git describe 2>/dev/null || echo "0.0.0")
RELEASE=release
RELEASE_FORMAT=${BINARY}-${GOOS}-${GOARCH}-${VERSION}.tar.gz

default: bin
.PHONY: bin
bin: deps ${BINARY}

.PHONY: release
release: release/gprobe-linux-amd64-${VERSION}.tar.gz
release: release/gprobe-linux-386-${VERSION}.tar.gz
release: release/gprobe-darwin-amd64-${VERSION}.tar.gz

.PHONY: release-dir
release-dir:
	mkdir -p ${RELEASE}

${RELEASE}/%-${VERSION}.tar.gz: | release-dir
	GOOS=$(shell echo $* | cut -d '-' -f 2) GOARCH=$(shell echo $* | cut -d '-' -f 3) \
	go build -ldflags="-s -w -X main.version=${VERSION} -v" -o ${RELEASE}/${BINARY}
	tar -C ${RELEASE} -czf $@ ${BINARY}
	rm ${RELEASE}/${BINARY}

${BINARY}:
	go build -ldflags="-s -w -X main.version=${VERSION} -v" -o $@

.PHONY: lint
lint:
	go get -u github.com/golang/lint/golint
	golint -set_exit_status ./...

.PHONY: test
test:
	go test -v $(go list ./... | grep -v /acctest/)

.PHONY: acctest
acctest: ${BINARY}
	go test -v ./acctest/... -args -gprobe `pwd`/${BINARY}

.PHONY: deps
deps:
	go get -u ./...

.PHONY: test-deps
test-deps:
	go get -t -u ./...

.PHONY: clean
clean:
	rm -f ${BINARY}
	rm -rf ${RELEASE}
