
FROM alpine AS builder

FROM scratch
ARG VERSION
ENV VERSION=${VERSION}
COPY noah-mqtt /
COPY LICENSE /
COPY passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER gouser
ENTRYPOINT ["/noah-mqtt"]