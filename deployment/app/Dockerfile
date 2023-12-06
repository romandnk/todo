FROM golang:1.21.3 as build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/app ./cmd/app
RUN chmod a+x /app/bin/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/

COPY ./config/ ./config/

CMD ["./bin/app"]