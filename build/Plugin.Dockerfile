ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.18 AS build

COPY . /home/workspace
WORKDIR /home/workspace

ARG UPX_VERSION=3.96
ARG TARGETPLATFORM

RUN apk --update add curl
RUN case ${TARGETPLATFORM:-linux/amd64} in \
        "linux/arm64")   UPX_ARCH="arm64" ;; \
        *)               UPX_ARCH="amd64" ;; \
    esac; \
    curl --location --output upx_linux.tar.xz "https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${UPX_ARCH}_linux.tar.xz"; \
	tar -xJf upx_linux.tar.xz; \
	cp upx-${UPX_VERSION}-${UPX_ARCH}_linux/upx /bin/

RUN go build -o bin/ ./cmd/plugin
RUN upx -4 bin/plugin


FROM alpine:3.18.2

COPY --from=build /home/workspace/bin/plugin /bin/velatemplatetesterplugin

ENTRYPOINT ["/bin/velatemplatetesterplugin"]
