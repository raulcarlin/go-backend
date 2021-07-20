FROM golang:1.13-alpine AS build-env

WORKDIR /go-backend

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

FROM alpine AS run-app

RUN apk update && apk add --no-cache bash mysql-client

COPY --from=build-env /go-backend/bin/app /go-backend/
COPY --from=build-env /go-backend/bin/migrate /go-backend/
COPY --from=build-env /go-backend/migrations /go-backend/migrations

COPY --from=build-env /go-backend/docker/app/bin /usr/local/bin/go-backend/
RUN chmod +x /usr/local/bin/go-backend/*

EXPOSE 8080

CMD ["/go-backend/app"]