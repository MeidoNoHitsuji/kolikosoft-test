FROM golang:1.25.7-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/app

RUN apk --update --no-cache add make

COPY Makefile go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build

FROM scratch

COPY --from=builder /go/src/app/app /app

ENTRYPOINT ["/app"]
