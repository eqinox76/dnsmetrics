FROM golang:1.21

ENV CGO_ENABLED=0
WORKDIR /app
COPY *.go go.mod go.sum .
RUN go build

FROM gcr.io/distroless/static-debian12
LABEL org.opencontainers.image.source="https://github.com/eqinoc76/dnsmetric"
COPY --from=0 /app/dnsmetric /bin/
ENTRYPOINT ["/bin/dnsmetric"]