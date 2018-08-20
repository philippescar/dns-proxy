# iron/go is the alpine image with only ca-certificates added
FROM iron/go

WORKDIR /app

ADD main /app/

EXPOSE 53

ENTRYPOINT ["./main"]