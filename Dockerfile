FROM golang:alpine

MAINTAINER seka <s.seka134704@gmail.com>

ADD . /go/src/github.com/seka/bbs-sample
RUN go install github.com/seka/bbs-sample/cmd/...
WORKDIR /go/src/github.com/seka/bbs-sample

EXPOSE 8080

CMD ["bbs-sampled"]
