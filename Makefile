# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gbfe android ios gbfe-cross evm all test clean
.PHONY: gbfe-linux gbfe-linux-386 gbfe-linux-amd64 gbfe-linux-mips64 gbfe-linux-mips64le
.PHONY: gbfe-linux-arm gbfe-linux-arm-5 gbfe-linux-arm-6 gbfe-linux-arm-7 gbfe-linux-arm64
.PHONY: gbfe-darwin gbfe-darwin-386 gbfe-darwin-amd64
.PHONY: gbfe-windows gbfe-windows-386 gbfe-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gbfe:
	$(GORUN) build/ci.go install ./cmd/gbfe
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gbfe\" to launch gbfe."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gbfe.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gbfe-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gbfe.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gbfe-cross: gbfe-linux gbfe-darwin gbfe-windows gbfe-android gbfe-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-*

gbfe-linux: gbfe-linux-386 gbfe-linux-amd64 gbfe-linux-arm gbfe-linux-mips64 gbfe-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-*

gbfe-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gbfe
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep 386

gbfe-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gbfe
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep amd64

gbfe-linux-arm: gbfe-linux-arm-5 gbfe-linux-arm-6 gbfe-linux-arm-7 gbfe-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep arm

gbfe-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gbfe
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep arm-5

gbfe-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gbfe
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep arm-6

gbfe-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gbfe
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep arm-7

gbfe-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gbfe
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep arm64

gbfe-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gbfe
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep mips

gbfe-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gbfe
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep mipsle

gbfe-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gbfe
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep mips64

gbfe-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gbfe
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-linux-* | grep mips64le

gbfe-darwin: gbfe-darwin-386 gbfe-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-darwin-*

gbfe-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gbfe
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-darwin-* | grep 386

gbfe-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gbfe
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-darwin-* | grep amd64

gbfe-windows: gbfe-windows-386 gbfe-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-windows-*

gbfe-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gbfe
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-windows-* | grep 386

gbfe-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gbfe
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbfe-windows-* | grep amd64
