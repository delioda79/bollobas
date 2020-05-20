FROM golang:1.13 as builder
RUN cd ..
RUN mkdir bollobas
WORKDIR bollobas
COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o bollobas ./cmd/bollobas/main.go && wget https://github.com/golang-migrate/migrate/releases/download/v4.1.0/migrate.linux-amd64.tar.gz && tar -xvf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 migrate

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bollobas/bollobas .


CMD ["./bollobas"]
