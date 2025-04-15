FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go env -w GO111MODULE=on && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-job-master cmd/master/master.go


FROM alpine:latest

ARG APP_ROOT=/app/
ARG APP_NAME=go-job-master
ENV RUN=${APP_ROOT}${APP_NAME}
USER root

WORKDIR ${APP_ROOT}

COPY --from=builder ${APP_ROOT}${APP_NAME} ${APP_ROOT}

CMD ["/bin/sh", "-c", "$RUN"]

