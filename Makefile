clean:
	rm coverage.out || true
	rm coverage.html || true
	rm docker/velatemplatetesterapi || true
	rm docker/velatemplatetesterplugin || true
test-api:
	go test app_test.go validator_test.go util_test.go app.go validator.go util.go -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
test-plugin:
	go test plugin_test.go validator_test.go util_test.go plugin.go validator.go util.go -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html	
check:
	gofmt -l -w -s .
	go vet app.go validator.go util.go
	go vet plugin.go validator.go util.go
	go test app_test.go validator_test.go util_test.go app.go validator.go util.go
	go test plugin_test.go validator_test.go util_test.go plugin.go validator.go util.go
coveralls:
	go test plugin_test.go validator_test.go util_test.go plugin.go validator.go util.go -v -covermode=count -coverprofile=coverage.out -json > test-report.json
	go get github.com/mattn/goveralls
	${GOPATH}/bin/goveralls -coverprofile=coverage.out
run:
	go run app.go validator.go util.go || true
build-api:
	go build -o docker/velatemplatetesterapi app.go validator.go util.go
build-plugin:
	go build -o docker/velatemplatetesterplugin plugin.go validator.go util.go