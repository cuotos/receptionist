FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app

COPY . /app

RUN go build -o receptionist main.go

FROM alpine

EXPOSE 8080

COPY --from=builder /app/receptionist /receptionist

CMD ["/receptionist"]