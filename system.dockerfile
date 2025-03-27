FROM golang:1.23 
WORKDIR /src
ADD . /src
RUN go mod download
RUN go build -o /bin/system ./cmd/system/main.go
EXPOSE 3000 3434 4343 
RUN rm -rf /src 
CMD [ "/bin/system" ]