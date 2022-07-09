# build stage
FROM golang:1.18-alpine AS build

ADD . /go/build
WORKDIR /go/build
RUN go build -o website ./main.go

# final stage
FROM alpine:3.15

COPY --from=build /go/build/website /var/application/website
COPY --from=build /go/build/asset/html /var/application/asset/html
COPY --from=build /go/build/asset/html/portfolio /var/application/asset/html/portfolio
COPY --from=build /go/build/asset/template /var/application/asset/template


RUN apk add openssl --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/main && \
    apk add librdkafka=1.1.0-r0 --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community && \
    apk add --no-cache ca-certificates && \
    apk add --no-cache tzdata && \
    apk del openssl && rm -f /var/cache/apk/*

RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /var/application/website

EXPOSE 443

WORKDIR /var/application
CMD [ "./website" ]