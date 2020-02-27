# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

ARG LOG_DIR=/app/logs

RUN mkdir -p ${LOG_DIR}

ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.cn

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

VOLUME ${LOG_DIR}

CMD ["./main"]