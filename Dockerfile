
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

RUN cd build && GOOS=linux go build -a -installsuffix cgo -o main main.go
RUN cp /go/src/build/main /main
EXPOSE 4000
ENV DB_FILE /db/db.sqlite3
CMD ["/main"]