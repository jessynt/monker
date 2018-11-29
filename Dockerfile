FROM golang:1.11-stretch as builder
ADD . /src
RUN cd /src && CGO_ENABLED=0 make all

FROM alpine as prod

COPY --from=builder /src/bin/monker .

VOLUME ["/etc/monker"]

CMD ["./monker", "-c", "/etc/monker/monker.yaml"]