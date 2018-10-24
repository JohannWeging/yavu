FROM golang:alpine AS builder

COPY . /go/src/github.com/JohannWeging/yavu
WORKDIR /go/src/github.com/JohannWeging/yavu

RUN set -x \
 && go install github.com/JohannWeging/yavu

FROM johannweging/base-alpine:latest

COPY --from=builder /go/bin/yavu /usr/bin

RUN set -x \
 && createuser yavu
 
ENTRYPOINT ["/usr/bin/dumb-init" "--"]
CMD ["yavu"]
