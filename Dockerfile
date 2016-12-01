FROM daocloud.io/jaggerwang/go

ENV APP_PATH=/go/src/jaggerwang.net/zqcserverdemo
ENV DATA_PATH=/data

ADD . $APP_PATH
WORKDIR $APP_PATH

RUN ./build.sh

VOLUME $DATA_PATH

EXPOSE 1323

CMD supervisord
