FROM alpine:3.4

ADD . /v1-metadata-publisher/

RUN apk add --update bash \
  && apk --update add git go ca-certificates \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/v1-metadata-publisher" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && mv v1-metadata-publisher/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go build \
  && go test ./... \
  && mv v1-metadata-publisher /v1-metadata-publisher-app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*

CMD [ "/v1-metadata-publisher-app" ]