FROM golang:alpine as builder
WORKDIR /go/src/eliranbarnoy/drone-cloudwatch-version-metric
COPY main.go  .
RUN apk add --update git
RUN go get ./.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
COPY --from=builder /go/src/eliranbarnoy/drone-cloudwatch-version-metric/app .
CMD ["./app"]  
