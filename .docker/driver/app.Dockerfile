FROM alpine

WORKDIR /driver

COPY --from=driver-build:develop /app/cmd/driver/driver ./driver

CMD ["/driver/driver"]
