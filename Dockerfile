FROM quay.io/projectquay/golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG CGO_ENABLED
RUN make build TARGETARCH=${TARGETARCH} TARGETOS=${TARGETOS} CGO_ENABLED=${CGO_ENABLED}


FROM scratch
WORKDIR /
COPY --from=builder /go/src/app/consumer .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./consumer", "start"]