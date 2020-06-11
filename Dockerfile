FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod /app

RUN go mod download

COPY . /app

RUN go build -o bin/receptionist .

FROM alpine

EXPOSE 8080

COPY --from=builder /app/bin/receptionist /receptionist

CMD ["/receptionist"]