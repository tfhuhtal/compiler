FROM golang:1.23-alpine

RUN apk add --no-cache bash binutils

WORKDIR /compiler
COPY go.mod .
RUN go mod tidy

COPY . .

RUN go test ./...

EXPOSE 3000

CMD ["go", "run", "main.go", "serve", "--host=0.0.0.0"]
