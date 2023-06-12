docker_tag=latest

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
	go get github.com/mattn/goveralls@v0.0.9
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
	go mod tidy
	go test -v ./... -tags test
	mkdir -p bin
	go build -o bin/ ./...
integration-test:
	go test -v ./... -tags integration
functional-test-plugin:
	docker pull devatherock/vela-template-tester:latest
	go test -v ./... -tags functional
functional-test-api:
	docker pull devatherock/vela-template-tester-api:$(docker_tag)
	DOCKER_TAG=$(docker_tag) docker-compose -f build/docker-compose.yml up -d
	sleep 1
	go test -v ./... -tags api
	docker-compose -f build/docker-compose.yml down