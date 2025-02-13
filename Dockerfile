FROM golang:1.23

ENV sourcesdir=/go/src/github.com/microservices-demo/user/
ENV MONGO_HOST=user-db:27017
ENV HATEAOS=user
ENV USER_DATABASE=users

COPY . ${sourcesdir}

RUN go install github.com/DataDog/orchestrion@latest

RUN cd /go/src/github.com/microservices-demo/user/ \
    && orchestrion pin \
    # && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app github.com/microservices-demo/user/
    && CGO_ENABLED=0 GOOS=linux go build -toolexec="orchestrion toolexec" -a -installsuffix cgo -o /app github.com/microservices-demo/user/

# RUN  update
# RUN apk add git
# RUN go get -v github.com/Masterminds/glide && cd ${sourcesdir} && glide install && go install

# ENTRYPOINT user
# EXPOSE 8084
CMD ["/app", "-port=8080"]
EXPOSE 8080
