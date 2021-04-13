# TODO
# change to smaller image like alpine

FROM ubuntu
COPY ./misc-go /misc-go
ENTRYPOINT /misc-go

