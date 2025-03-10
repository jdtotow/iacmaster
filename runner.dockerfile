FROM golang:1.23 as builder 
WORKDIR /src
ADD . /src
RUN go mod download
RUN go build -o ./bin/runner ./cmd/runner/main.go

FROM ubuntu:24.10
RUN apt-get update && apt-get install -y git curl unzip 
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv
RUN ln -s ~/.tfenv/bin/* /usr/local/bin
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN mkdir /app 
COPY --from=builder /src/bin/runner /app/runner
RUN mkdir /runner 
EXPOSE 8787 
CMD ["/app/runner"]
