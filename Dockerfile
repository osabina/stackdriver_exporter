FROM us.gcr.io/rsg-base-prod/golang:1.17.11 AS builder

ARG ARCH="amd64"
ARG OS="linux"

COPY . /src

RUN cd /src && go build -mod=vendor -o /tmp/stackdriver_exporter

FROM gcr.io/distroless/base

COPY --from=builder /tmp/stackdriver_exporter /usr/bin/

ENTRYPOINT ["/bin/stackdriver_exporter"]