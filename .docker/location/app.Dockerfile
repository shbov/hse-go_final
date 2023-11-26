FROM alpine

WORKDIR /app

COPY --from=build:latest /app/cmd/location/app ./app

CMD ["/app/app", "-c", "config.yaml"]
