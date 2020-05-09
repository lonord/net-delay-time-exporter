FROM --platform=$BUILDPLATFORM golang:1.13-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

COPY . /project
RUN cd /project \
    && export GOOS=$(echo $TARGETPLATFORM | cut -d "/" -f 1) \
    && export GOARCH=$(echo $TARGETPLATFORM | cut -d "/" -f 2) \
    && export GOARM=$(echo $TARGETPLATFORM | cut -d "/" -f 3) \
    && export GOARM=${GOARM#v} \
    && go build -o /net-delay-time-exporter -ldflags "-s -w" .

FROM alpine:3.11

ENV SERVERS "github.com"
ENV LISTEN ":8080"

COPY --from=builder /net-delay-time-exporter /net-delay-time-exporter
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]