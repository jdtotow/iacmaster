FROM golang:1.19
WORKDIR /app
COPY ./bin/system /
EXPOSE 3000
# Run
CMD ["/system"]