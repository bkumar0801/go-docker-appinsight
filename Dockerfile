FROM golang:1.9.4 as builder

# Set the working directory to the app directory
WORKDIR /go/src/go-docker-appinsight

# Install godeps
RUN go get -d github.com/Microsoft/ApplicationInsights-Go/appinsights

# Copy the application files
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hello-world main.go


FROM alpine
 
RUN apk --no-cache add ca-certificates

WORKDIR /app
 
COPY --from=builder /go/src/go-docker-appinsight/hello-world ./

ENTRYPOINT /app/hello-world

EXPOSE 8080
