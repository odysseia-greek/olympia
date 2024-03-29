FROM golang:1.22-alpine

# Install oapi-codegen
RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
WORKDIR /models
CMD ["tail", "-f", "/dev/null"]
