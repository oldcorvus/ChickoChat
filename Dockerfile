FROM golang:1.20-buster



WORKDIR /go/src/app
COPY . .


# Build the Go app
RUN GOPROXY=https://goproxy.cn go get github.com/klauspost/compress
RUN go build -o chat .

# Expose port 8080 to the outside world
EXPOSE 8009

# Command to run the executable
CMD ["./chat"]