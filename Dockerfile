FROM cr.yandex/secured/golang as builder
ARG devArgs=""

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 go build -gcflags "${devArgs}" -mod=vendor -o bin/service ./cmd/main.go

RUN mkdir -p /certs/postgres && \
    wget -O /certs/postgres/postgre-cluster-tls.crt https://storage.yandexcloud.net/secured/postgre-cluster-tls.crt && \
    chmod 0600 /certs/postgres

FROM cr.yandex/secured/debian

COPY --from=builder /app/bin/service /service
COPY --from=builder /certs /certs

ENV GOTRACEBACK=single

COPY migrations /migrations
COPY testdata /testdata
CMD ["/service"]
