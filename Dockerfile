ARG EXECUTABLE
FROM dep AS builder
ARG SERVICE
ARG EXECUTABLE

# Build the application.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o /go/bin/${EXECUTABLE} /go/src/app/${SERVICE}/cmd/${EXECUTABLE}

ENTRYPOINT /go/bin/${EXECUTABLE}

FROM alpine:latest
ARG EXECUTABLE
ENV EXECUTABLE=${EXECUTABLE}
# Copy executable to /bin/
COPY --from=builder /go/bin/${EXECUTABLE} /bin/${EXECUTABLE}

ENTRYPOINT /bin/${EXECUTABLE}