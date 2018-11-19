FROM golang:1.11 AS builder

# set working directory
WORKDIR /build

# copy our entire structure
COPY . ./

# build the app (statically)
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /portal .

# create the image from scratch
FROM scratch

# copy the app
COPY --from=builder /portal /usr/bin/local/portal

# copy the static files
COPY --from=builder /build/static/ /var/www/static/

# copy root certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# set entrypoint (our application)
ENTRYPOINT ["/usr/bin/local/portal"]
