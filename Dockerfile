# build stage
FROM golang:1.18-alpine AS build

ADD . /go/build
WORKDIR /go/build
RUN go build -o website ./main.go

# final stage
FROM golang:alpine3.15

COPY --from=build /go/build/website /var/application/website
COPY --from=build /go/build/asset/html /var/application/asset/html
COPY --from=build /go/build/asset/html/portfolio /var/application/asset/html/portfolio
COPY --from=build /go/build/asset/template /var/application/asset/template

RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /var/application/website

EXPOSE 80

WORKDIR /var/application
CMD [ "./website" ]