FROM golang:1.11 AS builder

# set working directory
WORKDIR /build

# copy our entire structure
COPY . ./

# build the app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /portal .

# create our actual image with the compiled binary
FROM scratch
COPY --from=builder /portal /usr/bin/local/portal
ENTRYPOINT ["/usr/bin/local/portal"]
