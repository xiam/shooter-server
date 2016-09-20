FROM centos:7

RUN mkdir /app

COPY shooter_linux_amd64 /app/shooter-server

WORKDIR /app

ENTRYPOINT [ "/app/shooter-server", "-listen", "0.0.0.0:3223"]
