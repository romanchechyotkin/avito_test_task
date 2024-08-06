FROM golang:1.22-alpine AS builder

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin cmd/app/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/bin .

CMD [ "/app/bin" ]