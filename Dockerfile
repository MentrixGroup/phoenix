FROM golang:alpine
WORKDIR /app  
COPY .  .
RUN go version
RUN ls -la
RUN apk add --no-cache make
RUN cd service/ && make build && \
    cd ../import/ && make build && \
    cd ../event-bridge/stream/ && make  && \
    cd ../publish/ && make 
