FROM golang:1.23 AS base

ENV sourcesdir=/go/src/github.com/microservices-demo/user/
COPY . ${sourcesdir}

RUN go install github.com/DataDog/orchestrion@latest

RUN cd /go/src/github.com/microservices-demo/user/ \
    && orchestrion pin \
    # && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app github.com/microservices-demo/user/
    && CGO_ENABLED=0 GOOS=linux go build -toolexec="orchestrion toolexec" -a -installsuffix cgo -o /user github.com/microservices-demo/user/

FROM golang:1.23

WORKDIR / 

ENV MONGO_HOST=user-db:27017 \
    HATEAOS=user \
    USER_DATABASE=mongodb

EXPOSE 80
COPY --from=base /user /

CMD ["/user", "-port=80"]
