FROM golang:1.23-alpine AS builder
WORKDIR /src
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=off
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/opspulse-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
ENV OPS_SERVER_ADDR=:8080
ENV OPS_DB_PATH=/data/opspulse.db
COPY --from=builder /out/opspulse-server /usr/local/bin/opspulse-server
VOLUME ["/data"]
EXPOSE 8080
CMD ["opspulse-server"]
