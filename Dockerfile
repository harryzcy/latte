FROM golang:1.21.4-alpine3.18 as builder

WORKDIR /app

COPY . ./

RUN set -ex && \
  go mod download && \
  go build \
  -ldflags="-w -s" \
  -o /bin/snuuze

FROM alpine:3.18.4

COPY --from=builder /bin/snuuze /bin/snuuze

EXPOSE 1323

CMD ["/bin/snuuze"]
