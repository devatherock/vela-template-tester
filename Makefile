docker_tag=latest
skip_pull=false

clean:
	rm -f coverage.out
	rm -f coverage.html
	rm -f test-report.json
	rm -rf bin
	go clean -testcache
check:
	gofmt -l -w -s .
	go vet ./...
	go test -v ./... -tags test -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
coveralls:
	go test -v ./... -tags test -covermode=count -coverprofile=coverage.out -json > test-report.json
	go get github.com/mattn/goveralls
	${GOPATH}/bin/goveralls -coverprofile=coverage.out
run-api:
	go build -o bin/ ./...
	./bin/app
run-plugin:
	go build -o bin/ ./...
	./bin/plugin	
build-all:
	gofmt -l -w -s .
	go vet ./...
	go test -v ./... -tags test
	mkdir -p bin
	go build -o bin/ ./...
integration-test:
	go test -v ./... -tags integration
docker-build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/ ./...
	docker build -t devatherock/vela-template-tester:$(docker_tag) \
	    -f build/Plugin.Dockerfile .
functional-test-plugin:
ifneq ($(skip_pull), true)
	docker pull devatherock/vela-template-tester:$(docker_tag)
endif
	go test -v ./... -tags functional
functional-test-api:
	docker pull devatherock/vela-template-tester-api:$(docker_tag)
	DOCKER_TAG=$(docker_tag) docker-compose -f build/docker-compose.yml up -d
	sleep 1
	go test -v ./... -tags api
	docker-compose -f build/docker-compose.yml down