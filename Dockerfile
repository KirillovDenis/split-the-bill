FROM golang:1.16

WORKDIR /app/split-the-bill

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o ./out/split-the-bill -v ./...

CMD ["./out/split-the-bill"]