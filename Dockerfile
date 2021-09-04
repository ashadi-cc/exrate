## COMPILE API APP
FROM golang:1.16-alpine as build

#update apk
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app
COPY . .
#build api app
RUN go build -o build/api cmd/api/main.go

## CREATE BUILD CONTAINER FROM alpine to REDUCE IMAGE SIZE
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /app/build/api /app/api
WORKDIR /app

EXPOSE 8080
CMD [ "./api" ]