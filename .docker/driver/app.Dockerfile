FROM alpine

WORKDIR /app

COPY --from=build:latest /app/cmd/driver/app ./app

CMD ["/app/app", "-c", "config.yaml"]
