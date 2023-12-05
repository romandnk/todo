FROM golang:1.21.4 as build

WORKDIR /app

COPY . .
RUN go mod download

ARG BUILD_VERSION

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/app \
    -ldflags "-X main.release=$BUILD_VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(git log --format="%h" -n 1)" \
    ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/

COPY ./config/ ./config/

CMD ["./bin/app"]