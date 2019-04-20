# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.10
ARG GO_VERSION=1.12

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git

# Set the environment variables for the go command:
# * CGO_ENABLED=0 to build a statically-linked executable
# * GOFLAGS=-mod=vendor to force `go build` to look into the `/vendor` folder.
#ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

# Import the code from the context.
COPY ./ /app

# Build the executable to `/app`. Mark the build as statically linked.
RUN go build \
    -installsuffix 'static' \
    -o ./bin/app .

# Final stage: the running container.
FROM scratch AS final

WORKDIR /app

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the second stage.
COPY --from=builder /app/bin/app /app

# Perform any further action as an unprivileged user.
USER nobody:nobody

## Run the compiled binary.
ENTRYPOINT ["./app"]