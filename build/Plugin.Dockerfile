ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.18 AS build

COPY . /home/workspace
WORKDIR /home/workspace

ARG TARGETPLATFORM
RUN case ${TARGETPLATFORM:-linux/amd64} in \
        "linux/arm64")   GO_ARCH="arm64" ;; \
        *)               GO_ARCH="amd64" ;; \
    esac; \
    GOARCH=${GO_ARCH} go build -o bin/ ./cmd/plugin


FROM alpine:3.18.0

COPY --from=build /home/workspace/bin/plugin /bin/velatemplatetesterplugin

ENTRYPOINT ["/bin/velatemplatetesterplugin"]
