FROM golang:alpine AS builder

COPY . /go/src/github.com/yoo/yavu
WORKDIR /go/src/github.com/yoo/yavu

RUN set -x \
 && go install -mod=vendor github.com/yoo/yavu

FROM alpine:latest

COPY --from=builder /go/bin/yavu /usr/bin

RUN set -x \
 && apk add --update --no-cache dumb-init ca-certificates \
 && adduser -D yavu
 
USER yavu

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/yavu"]
