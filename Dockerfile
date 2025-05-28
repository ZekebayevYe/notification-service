# syntax=docker/dockerfile:1
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o notification ./cmd/notify

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/notification .

ENTRYPOINT ["/notification"]
