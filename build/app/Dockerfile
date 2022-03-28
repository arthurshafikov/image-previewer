FROM golang:1.17.3 as build

ENV BIN_FILE /opt/app/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS

RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/main.go

FROM alpine:3.9

LABEL MAINTAINERS="thewolf27"

ENV BIN_FILE "/opt/app/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CODE_DIR /go/src/

ENV CONFIG_FOLDER "/opt/app/configs"
COPY --from=build ${CODE_DIR}/configs ${CONFIG_FOLDER}

ENV STORAGE_FOLDER "/opt/app/storage"
COPY --from=build ${CODE_DIR}/storage ${STORAGE_FOLDER}

CMD ${BIN_FILE} --configFolder ${CONFIG_FOLDER} --storageFolder ${STORAGE_FOLDER}
