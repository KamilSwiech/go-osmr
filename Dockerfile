FROM golang:1.9 as builder  
RUN mkdir -p /go/src/go-osmr   
WORKDIR /go/src/go-osmr    
COPY src/ .  
RUN go get -d  
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go-osmr .    

FROM alpine:latest  
WORKDIR /
COPY --from=builder /go/src/folder/go-osmr .  
CMD ["/go-osmr"]
