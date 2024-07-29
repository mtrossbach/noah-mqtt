
FROM registry.access.redhat.com/ubi8/ubi-minimal AS builder
RUN microdnf install -y shadow-utils && microdnf clean all
RUN set -o pipefail && groupadd -r -g 1000 gouser && useradd -r -u 1000 -g gouser -m -d /opt/gouser -s /bin/bash gouser && cat /etc/passwd | grep gouser > /etc/passwd_gouser

FROM scratch
ARG VERSION
ENV VERSION=${VERSION}
COPY noah-mqtt /
COPY LICENSE /
COPY --from=builder /etc/pki/tls/certs/ca-bundle.crt /etc/pki/tls/certs/
COPY --from=builder /etc/passwd_gouser /etc/passwd
USER gouser
ENTRYPOINT ["/noah-mqtt"]