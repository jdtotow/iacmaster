FROM golang:1.23 
WORKDIR /src
ADD . /src
RUN go mod download
RUN go build -o /bin/proxy ./cmd/proxy/main.go
EXPOSE 5454 
RUN rm -rf /src 
CMD [ "/bin/proxy" ]