FROM golang:1.22-alpine3.18 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserverx .

FROM scratch
COPY --from=builder ["/build/apiserverx", "/"]
EXPOSE 80
ENTRYPOINT ["/apiserverx"]
