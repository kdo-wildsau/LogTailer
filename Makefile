DATE = $(shell date +%FT%T%Z)
GIT_VER=$(shell git rev-parse HEAD)

BUILD_DIR = build
APP = logTailer
CMD = cmd/logTailer/main.go
MC_DIR = build/bin
LOG_DIR= build/log

LDFLAGS=-ldflags "-X main.version=${GIT_VER}"

.PHONY: dut monitor.darwin64 monitor.darwinArm lint clean distclean mrproper

GO_TOOLS := gridx/golang-dev:1.12.9.latest-linux-amd64
GO_PROJECT := github.com/kdo-wildsau/logTailer

# Build the project
all:
	@echo "cmd:"
	@echo "  dut	            build dut for arm"
	@echo ""
	@echo "  app            build all monitor for all os"
	@echo "  app.darwin64   build app for osx x64"
	@echo "  app.darwinArm  build app for osx arm64"
	@echo "  app.linux64    build app for linux arm64"
	@echo "  app.windows64  build app for windows x64"
	@echo ""
	@echo "  lint           run linters"
	@echo "  clean          remove dut binarys"
	@echo "  distclean      remove build folder"


app: app.darwin64 app.darwinArm app.linux64 app.windows64

app.darwin64:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-darwin -v ${CMD}

app.darwinArm:
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-darwin-arm -v ${CMD}

app.linux64:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-linux -v ${CMD}

app.windows64:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}_64.exe -v ${CMD}

lint:
	golint -set_exit_status $(shell go list ./...)

clean:
	-rm -f ${BUILD_DIR}/${DUT}-*

distclean:
	rm -rf ./build
