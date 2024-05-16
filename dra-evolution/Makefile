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

.PHONY: all
all: fmt test build

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	cd cmd/schedule && go build
	cd cmd/gen && go build
	cd cmd/mock-apiserver && go build
