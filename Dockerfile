FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod /app

RUN go mod download

COPY . /app

RUN echo "git rev-parse --short HEAD"

RUN go get github.com/markbates/pkger/cmd/pkger && \
  pkger list && \
  pkger && \
  GIT_COMMIT=$(git rev-parse --short HEAD) && \
  go build -ldflags "-X main.appVersion=git-$GIT_COMMIT" -o bin/receptionist .

FROM alpine

EXPOSE 8080

COPY --from=builder /app/bin/receptionist /receptionist

CMD ["/receptionist"]
