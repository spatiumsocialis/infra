ARG EXECUTABLE
ARG GOOGLE_GCR_HOSTNAME
ARG GOOGLE_PROJECT_ID
FROM ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest AS sourcer
ARG SERVICE
ARG EXECUTABLE

WORKDIR /go/src/app

# Copy service's code into container
COPY ./pkg/services/${SERVICE} ./pkg/services/${SERVICE}
COPY ./configs/services/${SERVICE} ./configs/services/${SERVICE}
COPY ./cmd/services/${SERVICE} ./cmd/services/${SERVICE}

# # Download and install imports
RUN go get -v ./...

FROM sourcer AS builder
ARG SERVICE
ARG EXECUTABLE

COPY --from=sourcer /go/src/app /go/src/app

# Build the application.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o /go/bin/${EXECUTABLE} /go/src/app/cmd/services/${SERVICE}/${EXECUTABLE}

ENTRYPOINT /go/bin/${EXECUTABLE}

FROM alpine:latest
ARG EXECUTABLE
ENV EXECUTABLE=${EXECUTABLE}
# Copy executable to /bin/
COPY --from=builder /go/bin/${EXECUTABLE} /bin/${EXECUTABLE}

ENTRYPOINT /bin/${EXECUTABLE}
