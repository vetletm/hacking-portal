FROM golang:1.11 AS builder

# set working directory
WORKDIR /build

# copy our entire structure
COPY . ./

# build the app
RUN GO111MODULE=on go build -o /portal .

# create our actual image with the compiled binary
FROM scratch
COPY --from=builder /portal /usr/bin/local/portal
ENTRYPOINT ["/usr/bin/local/portal"]
