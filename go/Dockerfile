FROM golang:latest
 
RUN mkdir -p /app
 
WORKDIR /app
 
ADD . /app
 
RUN go build ./pipelines.go
 
CMD ["./pipelines"]