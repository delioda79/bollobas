FROM golang:1.12-alpine3.10 as builder
RUN cd ..
RUN mkdir bollobas
WORKDIR bollobas
COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o bollobas ./cmd/bollobas/main.go 

FROM scratch
COPY --from=builder /go/bollobas/bollobas .
CMD ["./bollobas"]
