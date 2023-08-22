FROM golang:1.20.7-alpine3.18 as builder
RUN apk add --update make && apk add --update openssl
RUN apk add --update bash
RUN bash --version
RUN bash
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
ARG VERSION
ENV VERSION=$VERSION
RUN echo $VERSION 
WORKDIR /app
COPY . .
RUN make

FROM alpine:latest

# Create a group and user
RUN addgroup --gid 9999 hyper && adduser --disabled-password --gecos '' --uid 9999 -G hyper -s /bin/ash hyper

WORKDIR /app
COPY --from=builder /go/bin/hyper /app/.scripts/entrypoint-hyper.sh /app/

ENTRYPOINT ["/app/hyper"]
