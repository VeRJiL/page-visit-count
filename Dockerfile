FROM golang:1.20-rc-alpine as build-env

# Set environment variable
ENV APP_NAME page-visit-count
ENV CMD_PATH main.go

WORKDIR $GOPATH/src/$APP_NAME
COPY . $GOPATH/src/$APP_NAME

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH



# Build Stage
FROM alpine:3.17

# Set environment variable
ENV APP_NAME page-visit-count

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
COPY templates ./templates
RUN chown -R 777 /templates

# Expose application port
EXPOSE 8081

# Start app
CMD ./$APP_NAME
