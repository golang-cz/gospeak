# -----------------------------------------------------------------
# Builder
# -----------------------------------------------------------------
FROM golang:1.19-alpine3.16 as builder
ARG VERSION

RUN apk add --update git

ADD ./ /src

WORKDIR /src
RUN go build -ldflags="-s -w -X main.VERSION=${VERSION}" -o /usr/bin/go2webrpc ./cmd/go2webrpc

# -----------------------------------------------------------------
# Runner
# -----------------------------------------------------------------
FROM alpine:3.16

ENV TZ=UTC

RUN apk add --no-cache --update ca-certificates

COPY --from=builder /usr/bin/go2webrpc /usr/bin/

ENTRYPOINT ["/usr/bin/go2webrpc"]
