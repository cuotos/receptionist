FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app

COPY . /app

RUN go build -o receptionist main.go

FROM alpine

COPY --from=builder /app/receptionist /receptionist

CMD ["/receptionist"]