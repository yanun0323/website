# build stage
FROM golang:1.18-alpine AS build

ADD . /go/build
WORKDIR /go/build
RUN go build -o website ./main.go

# final stage
FROM alpine:3.15

COPY --from=build /go/build/website /var/application/website
COPY --from=build /go/build/internal/asset/html /var/application/internal/asset/html
COPY --from=build /go/build/internal/asset/html/portfolio /var/application/internal/asset/html/portfolio
COPY --from=build /go/build/internal/asset/template /var/application/internal/asset/template

EXPOSE 8080

WORKDIR /var/application
CMD [ "./website" ]