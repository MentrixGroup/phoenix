FROM golang:alpine
WORKDIR /app  
COPY .  .
RUN ls -la
RUN cd service/ && make build && \
    cd ../import/ && make build && \
    cd ../event-bridge/stream/ && make build && \
    cd ../publish/ && make build
