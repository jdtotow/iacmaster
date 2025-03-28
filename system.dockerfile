FROM golang:1.23 as builder
WORKDIR /src
ADD . /src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/system ./cmd/system/main.go 

### 
## Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=builder /src/bin/system /
EXPOSE 3000 3434 4343
# Set the entry point for the container
ENTRYPOINT ["/system"]

