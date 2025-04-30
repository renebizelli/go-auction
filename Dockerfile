FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o app ./cmd/auction/main.go 

FROM scratch
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/cmd/auction/.env .

EXPOSE 8080

CMD [ "./app" ]
