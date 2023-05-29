FROM --platform=linux/amd64 golang:1.20-alpine as build

WORKDIR /app

COPY . .

RUN go mod download

# generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./router/routes.go

# build application binary
RUN go build -o /bin/entrypoint .

FROM --platform=linux/amd64 alpine:latest 

COPY --from=build /bin/entrypoint /bin/entrypoint

CMD ["/bin/entrypoint"]

