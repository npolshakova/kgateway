ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output
DEPSGOBIN := $(OUTPUT_DIR)/.bin

# proto compiler installation
PROTOC_VERSION:=3.6.1
PROTOC_URL:=https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}
.PHONY: install-protoc
.SILENT: install-protoc
install-protoc:
	mkdir -p $(DEPSGOBIN)
	if [ $(shell ${DEPSGOBIN}/protoc --version | grep -c ${PROTOC_VERSION}) -ne 0 ]; then \
		echo expected protoc version ${PROTOC_VERSION} already installed ;\
	else \
		if [ "$(shell uname)" = "Darwin" ]; then \
			echo "downloading protoc for osx" ;\
			wget $(PROTOC_URL)-osx-x86_64.zip -O $(DEPSGOBIN)/protoc-${PROTOC_VERSION}.zip ;\
		elif [ "$(shell uname -m)" = "aarch64" ]; then \
			echo "downloading protoc for linux aarch64" ;\
			wget $(PROTOC_URL)-linux-aarch_64.zip -O $(DEPSGOBIN)/protoc-${PROTOC_VERSION}.zip ;\
		else \
			echo "downloading protoc for linux x86-64" ;\
			wget $(PROTOC_URL)-linux-x86_64.zip -O $(DEPSGOBIN)/protoc-${PROTOC_VERSION}.zip ;\
		fi ;\
		unzip $(DEPSGOBIN)/protoc-${PROTOC_VERSION}.zip -d $(DEPSGOBIN)/protoc-${PROTOC_VERSION} ;\
		mv $(DEPSGOBIN)/protoc-${PROTOC_VERSION}/bin/protoc $(DEPSGOBIN)/protoc ;\
		chmod +x $(DEPSGOBIN)/protoc ;\
		rm -rf $(DEPSGOBIN)/protoc-${PROTOC_VERSION} $(DEPSGOBIN)/protoc-${PROTOC_VERSION}.zip ;\
	fi

.PHONY: generate-agentgateway
generate-agentgateway: install-protoc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	$(DEPSGOBIN)/protoc --proto_path=. --go_out=. --go_opt=paths=source_relative listener.proto
	$(DEPSGOBIN)/protoc --proto_path=. --go_out=. --go_opt=paths=source_relative a2a/a2a_target.proto
	$(DEPSGOBIN)/protoc --proto_path=. --go_out=. --go_opt=paths=source_relative mcp/mcp_target.proto

