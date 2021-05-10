FROM golang:alpine
WORKDIR /app  
COPY .  .
RUN touch .config.mk
RUN go version
RUN ls -la
RUN apk add --no-cache make
RUN cd service/ && make  
RUN    cd ../import/ && make   
RUN    cd ../event-bridge/stream/ && make  
RUN    cd ../publish/ && make 
