FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -ldflags="-w -s" \
    -o /aws-s3-exporter \
    ./cmd/aws-s3-exporter

FROM gcr.io/distroless/static-debian12

COPY --from=builder /aws-s3-exporter /bin/aws-s3-exporter

EXPOSE 2112

ENTRYPOINT ["/bin/aws-s3-exporter"]

CMD ["--config", "/etc/aws-s3-exporter/config.yaml"]

