ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.18 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./...
RUN go test -v ./... -tags integration


FROM alpine:3.18.0

COPY --from=build /home/workspace/bin/plugin /bin/velatemplatetesterplugin

ENTRYPOINT ["/bin/velatemplatetesterplugin"]
