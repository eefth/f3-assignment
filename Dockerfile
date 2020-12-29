FROM golang:1.15.6-alpine3.12

RUN set -ex; \
    apk update; \ 
    apk add --no-cache git; \
    go get github.com/google/uuid 

RUN mkdir -p /go/src/github.com/eefth/f3-assignment
ADD ./app /go/src/github.com/eefth/f3-assignment/app

#RUN mkdir -p go/src/github.com/eefth/client
ADD ./client /go/src/github.com/eefth/f3-assignment/client

WORKDIR /go/src/github.com/eefth/f3-assignment/client/

CMD CGO_ENABLED=0 go test -v *go