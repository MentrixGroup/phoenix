FROM golang:alpine
WORKDIR /app  
COPY .  .
RUN go version
RUN ls -la
RUN apk add --no-cache make
RUN cd service/ && make build  
RUN    cd ../import/ && make build  
RUN    cd ../event-bridge/stream/ && make  
RUN    cd ../publish/ && make 
