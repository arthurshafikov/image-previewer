FROM golang:1.18 as build

ENV CODE_DIR /go/src/
ENV TEST_DIR ./tests

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

CMD go test ${TEST_DIR} -v
