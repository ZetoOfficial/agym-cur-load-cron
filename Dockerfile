FROM golang:1.22.3 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/main .
COPY config/config.yaml /root/config/config.yaml

CMD ["./main"]
