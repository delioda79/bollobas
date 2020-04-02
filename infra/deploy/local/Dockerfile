FROM golang:1.13 as builder
RUN cd ..
RUN mkdir bollobas
WORKDIR bollobas
COPY . ./
ARG version=dev
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o bollobas ./cmd/bollobas/main.go

FROM scratch
COPY --from=builder /go/bollobas/bollobas .
CMD ["./bollobas"]