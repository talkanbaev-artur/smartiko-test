FROM golang:alpine as build 
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go build -o /main src/cmd/*.go

FROM alpine:latest
WORKDIR /
COPY --from=build /main /
COPY migrations /migrations
ENTRYPOINT [ "/main" ]