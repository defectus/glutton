BINARY = glutton
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH ?= amd64
GOOS ?= linux


VERSION?=0.1.0
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags --abbrev=0)
BUILDTIME=$(shell date '+%Y.%m.%d-%H.%M.%S')

# Symlink into GOPATH
GITHUB_USERNAME=defectus
BUILD_DIR=.
CURRENT_DIR=$(shell pwd)
STATICS_DIR=definition

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w -X main.AUTHOR=${GITHUB_USERNAME} -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.TAG=${TAG} -X main.BUILDTIME=${BUILDTIME}"

# Build the project
all: clean test vet build

run: all
	./${BINARY}-${GOOS}-${GOARCH}

build: 
	cd ${BUILD_DIR}; \
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-${GOOS}-${GOARCH} cmd/main.go ; \
	cd - >/dev/null

test:
	go get github.com/tebeka/go2xunit
	cd ${BUILD_DIR}; \
	go test -v ./... 2>&1 | ${GOPATH}/bin/go2xunit -output ${TEST_REPORT} ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

docker:
	docker build -f ops/Dockerfile -t defectus/glutton -t defectus/glutton:${TAG} .
	docker push defectus/glutton defectus/glutton:${TAG}

.PHONY: test vet fmt clean docker