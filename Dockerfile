# Build server
FROM golang as build
ENV GO111MODULE=on
WORKDIR /seer

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/seer

# Update certificates
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# Aggregate into final image
FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/seer /bin/seer
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/bin/seer"]
