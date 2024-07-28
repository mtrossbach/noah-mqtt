FROM golang:alpine
ARG VERSION
ENV VERSION=${VERSION}
ENTRYPOINT ["/noah-mqtt"]
COPY noah-mqtt /
COPY LICENSE /
