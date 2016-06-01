FROM ubuntu
MAINTAINER = Gennadiy Dubina <gdubina@dataart.com>

RUN apt-get update && \
    apt-get install -y curl ipcalc vim

#test app
ENV WORK_DIR /opt/ivrtest

COPY ./files $WORK_DIR

WORKDIR $WORK_DIR
CMD ./run.sh 
