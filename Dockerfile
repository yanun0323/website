# build stage
FROM golang:1.18-alpine AS build

ADD . /go/build
WORKDIR /go/build
RUN go build -o website ./main.go

# final stage
FROM alpine:3.15

COPY --from=build /go/build/website /var/application/website
COPY --from=build /go/build/config /var/application/config
COPY --from=build /go/build/internal/resource /var/application/internal/resource

EXPOSE 8080

WORKDIR /var/application
CMD [ "./website" ]