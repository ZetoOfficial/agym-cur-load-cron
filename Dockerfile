FROM golang:1.22.3-alpine

RUN apk --no-cache add bash
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

CMD ["wait-for-it.sh", "postgres:5432", "--", "./main"]
