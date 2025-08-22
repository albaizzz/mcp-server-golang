FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server /app/server
EXPOSE 3001
CMD ["/app/server"]