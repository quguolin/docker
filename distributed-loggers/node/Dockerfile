# Start from the latest golang base image
FROM golang:1.12

# Set the Current Working Directory inside the container
WORKDIR /app

ENV GOPROXY="https://goproxy.io"

COPY go.mod  go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]