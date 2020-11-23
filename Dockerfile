# Build Image
FROM golang:1.15-buster as build


# General setup
ENV SERVICE_NAME=maze
ENV APP_DIR=/go/src/github.com/avanticaTest/${SERVICE_NAME}
WORKDIR ${APP_DIR}
COPY . .


#Go build
RUN CGO_ENABLED=0 GOOS=linux go build -x -a \
        -installsuffix cgo \
        github.com/avanticaTest/${SERVICE_NAME}/cmd/${SERVICE_NAME}


# Deployable Image
FROM alpine:3.10

ENV SERVICE_NAME=maze

# Add binary to bin directory
COPY --from=build /go/src/github.com/avanticaTest/${SERVICE_NAME}/${SERVICE_NAME} /


EXPOSE 8080