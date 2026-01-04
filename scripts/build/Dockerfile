FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY fs /app/fs
RUN chmod +x /app/fs

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./fs", "-c", "configs/fs.yml"]
# if you use the Configuration Center, fs.yml is changed to the Configuration Center configuration.
#CMD ["./fs", "-c", "configs/fs.yml", "-enable-cc"]
