
FROM golang:latest as builder
LABEL maintainer="Jean Ribes <jean@ribes.ovh>"

WORKDIR /go/src
RUN mkdir -p build \
  && go get -u github.com/golang/dep/... \
  && cd /go/src/github.com/golang/dep \
  && git checkout ${GODEP_VERSION} \
  && go install github.com/golang/dep/... \
  && mv /go/bin/dep /usr/bin

COPY Gopkg.lock build
COPY Gopkg.toml build

COPY main.go build

RUN cd build && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go # build

FROM alpine:3.11
COPY --from=builder /go/src/build/main .
EXPOSE 4000
ENV DB_FILE /db/db.sqlite3
CMD ["/main"]