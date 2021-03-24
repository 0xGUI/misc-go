FROM ubuntu
COPY ./misc-go /misc-go
ENTRYPOINT /misc-go

