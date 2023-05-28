clean:
	rm coverage.out || true
	rm coverage.html || true
	rm test-report.json || true
	rm docker/velatemplatetesterapi || true
	rm docker/velatemplatetesterplugin || true
test-api:
	go test -v -tags api -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
test-plugin:
	go test -v -tags plugin -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
check:
	gofmt -l -w -s .
	go vet -tags api
	go vet -tags plugin
	go test -tags api
	go test -tags plugin
coveralls:
	go test -v -tags plugin -covermode=count -coverprofile=coverage.out -json > test-report.json
	go get github.com/mattn/goveralls@v0.0.9
	${GOPATH}/bin/goveralls -coverprofile=coverage.out
run-api:
	go build -o docker/velatemplatetesterapi -tags api
	./docker/velatemplatetesterapi
run-plugin:
	go build -o docker/velatemplatetesterplugin -tags plugin
	./docker/velatemplatetesterplugin	
build-api:
	go build -o docker/velatemplatetesterapi -tags api
build-plugin:
	go build -o docker/velatemplatetesterplugin -tags plugin
integration-test:
	go test -v -tags integration