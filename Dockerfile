
FROM alpine AS builder

FROM scratch
COPY noah-mqtt /
COPY LICENSE /
COPY passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER gouser
ENTRYPOINT ["/noah-mqtt"]