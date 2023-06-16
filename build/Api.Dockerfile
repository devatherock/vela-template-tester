ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine3.18 AS build

COPY . /home/workspace
WORKDIR /home/workspace

RUN go build -o bin/ ./cmd/app


FROM alpine:3.18.2

COPY --from=build /home/workspace/bin/app /bin/velatemplatetesterapi

ENTRYPOINT ["/bin/velatemplatetesterapi"]