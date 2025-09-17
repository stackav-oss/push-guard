# Copyright 2025 Stack AV Co.
# SPDX-License-Identifier: Apache-2.0
SHELL=/bin/bash

version = $(shell jq -j '.PushGuardVersion' < ./build.json)
ldflags = -X 'push-guard/config.PushGuardVersion="$(shell echo -n "${version}" | base64 -w 0)"' \
	-X 'push-guard/config.Disclaimer="$(shell jq -j '.Disclaimer' < ./build.json | base64 -w 0)"' \
	-X 'push-guard/config.LogCollectorURL="$(shell jq -j '.LogCollectorURL' < ./build.json | base64 -w 0)"' \
	-X 'push-guard/config.ProtectedBranches="$(shell jq -j '.ProtectedBranches' < ./build.json | base64 -w 0)"' \
	-X 'push-guard/config.ProtocolAndDomainAllowList="$(shell jq -j '.ProtocolAndDomainAllowList' < ./build.json | base64 -w 0)"' \
	-X 'push-guard/config.DirectoryAllowList="$(shell jq -j '.DirectoryAllowList' < ./build.json | base64 -w 0)"' \
	-X 'push-guard/config.DirectoryRegexAllowList="$(shell jq -j '.DirectoryRegexAllowList' < ./build.json | base64 -w 0)"'

ifeq ("${os}","")
	os = $(shell tr '[:upper:]' '[:lower:]' <<< $(shell uname))
endif

ifeq ("${arch}","")
	arch = $(shell tr '[:upper:]' '[:lower:]' <<< $(shell uname -m) | sed 's/x86_/amd/')
endif

binary_path := "build/${version}/${os}/${arch}/push-guard"


.PHONY: build
build:
	@echo "Building ${version}/${os}/${arch} ..."
	@echo "Installing dependencies ..."
	@go get -C . .
	@if [[ "${os}" = "windows" ]]; then \
	    echo "Building binary: \"${binary_path}.exe\""; \
	    env GOOS="${os}" GOARCH="${arch}" go build -C . -o "${binary_path}.exe" -ldflags "${ldflags}"; \
	else \
	    echo "Building binary: \"${binary_path}\""; \
	    env GOOS="${os}" GOARCH="${arch}" go build -C . -o "${binary_path}" -ldflags "${ldflags}"; \
	fi


clean:
	@[[ -d build ]] && { echo "Cleaning build ..."; rm -rf build; } || echo "build not found"
	@[[ -f go.sum ]] && { echo "Cleaning go.sum ..."; rm -f go.sum; } || echo "go.sum not found"
	@[[ -f junit.xml ]] && { echo "Cleaning junit.xml ..."; rm -f junit.xml; } || echo "junit.xml not found"


test:
	@echo "Installing gotestsum ..."
	@go install gotest.tools/gotestsum@latest
	@echo "Installing dependencies ..."
	@go get -C . .
	@echo "Generating test report ..."
	@gotestsum --junitfile junit.xml --format testname

