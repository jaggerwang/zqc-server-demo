FROM daocloud.io/jaggerwang/go

ENV APP_PATH=/go/src/zqc
ENV DATA_PATH=/data

ADD . $APP_PATH
WORKDIR $APP_PATH

RUN go get -d -v ./...
RUN go install -v .

VOLUME $DATA_PATH

EXPOSE 1323

CMD supervisord
