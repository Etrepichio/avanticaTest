# Build Image
FROM golang:1.15-buster as build


# General setup
ENV SERVICE_NAME=maze
ENV APP_DIR=/go/src/github.com/avanticaTest/${SERVICE_NAME}
WORKDIR ${APP_DIR}
COPY . .


# GO setup, build and test
RUN go get github.com/jteeuwen/go-bindata/... github.com/jstemmer/go-junit-report &&\
    go test ./... -v -coverprofile=cover.out | go-junit-report > report.xml &&\
    go tool cover -html=cover.out -o coverage_report.html

RUN CGO_ENABLED=0 GOOS=linux go build -x -a \
        -installsuffix cgo \
        github.com/avanticaTest/${SERVICE_NAME}/cmd/${SERVICE_NAME}

# Deployable Image
FROM alpine:3.10

ENV SERVICE_NAME=maze

# Add binary to bin directory
COPY --from=build /go/src/github.com/avanticaTest/${SERVICE_NAME}/${SERVICE_NAME} /
COPY --from=build /go/src/github.com/avanticaTest/${SERVICE_NAME}/coverage_report.html /
COPY --from=build /go/src/github.com/avanticaTest/${SERVICE_NAME}/report.xml /

ENV PORT=8080

# Command to execute at start up
ENTRYPOINT ["/maze"]

EXPOSE ${PORT}