FROM alpine

WORKDIR /location

COPY --from=location-build:develop /app/cmd/location/location ./location

CMD ["/location/location"]
