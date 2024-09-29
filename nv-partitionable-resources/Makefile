# Copyright 2024 The Kubernetes Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.DEFAULT_GOAL := all

DEEPCOPY_SOURCES += pkg/nvdevice
DEEPCOPY_SOURCES += pkg/resource

.PHONY: all
all: fmt test build

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	cd cmd/print-model && go build
	cd cmd/print-possible-allocations && go build

.PHONY: print-possible-allocations
print-model:
	go run ./cmd/print-model/...

.PHONY: print-possible-allocations
print-possible-allocations:
	go run ./cmd/print-possible-allocations/...

.PHONY: generate-deepcopy
generate-deepcopy:
	for dir in $(DEEPCOPY_SOURCES); do \
	    controller-gen \
	        object:headerFile=$(CURDIR)/hack/boilerplate.go.txt,year=$(shell date +"%Y") \
	        paths=$(CURDIR)/$${dir}/ \
	        output:object:dir=$(CURDIR)/$${dir}; \
	done

.PHONY: clean
clean:
	rm -f cmd/print-model/print-model
	rm -f cmd/print-possible-allocations/print-possible-allocations
	find . -name "zz_generated.*" | xargs rm -f
