
FROM alpine AS builder
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN set -o pipefail && addgroup --gid 1000 gouser && adduser --disabled-password --no-create-home --ingroup gouser gouser
RUN cat /etc/passwd
RUN cat /etc/passwd | grep gouser > /etc/passwd_gouser

FROM scratch
ARG VERSION
ENV VERSION=${VERSION}
COPY noah-mqtt /
COPY LICENSE /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd_gouser /etc/passwd
USER gouser
ENTRYPOINT ["/noah-mqtt"]