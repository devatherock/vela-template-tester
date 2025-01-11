ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-alpine3.20 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./cmd/plugin


FROM alpine:3.21.2

COPY --from=build /home/workspace/bin/plugin /bin/velatemplatetesterplugin

ENTRYPOINT ["/bin/velatemplatetesterplugin"]
