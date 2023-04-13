ARTIFACTS_DIR=artifacts/${VERSION}
GITHUB_USERNAME=sachaos

.PHONY: install
install: prepare
	go install

.PHONY: build
build: prepare
	go build

.PHONY: test
test: prepare
	go test -v

.PHONY: prepare
prepare: filter_parser.go

.PHONY: release
release: prepare
	GOOS=windows GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/todoist_windows_amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/todoist_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/todoist_darwin_arm64
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/todoist_linux_amd64
	ghr -u $(GITHUB_USERNAME) -t $(shell cat github_token) --replace ${VERSION} $(ARTIFACTS_DIR)

filter_parser.go: filter_parser.y
	go get golang.org/x/tools/cmd/goyacc
	goyacc -o filter_parser.go filter_parser.y
	rm y.output

docker-build:
	docker build -t todoist --build-arg TODOIST_API_TOKEN=$(token) .

docker-run:
	docker run -it todoist /bin/bash
