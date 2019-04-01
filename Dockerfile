FROM node as build-frontend
ADD . /easybot
WORKDIR /easybot/dashboard
RUN npm install
RUN CI=false npm run build

FROM golang AS build-env
ADD . /go/src/easybot
COPY --from=build-frontend /easybot/dashboard/build /go/src/easybot/dashboard/build
WORKDIR /go/src/easybot
RUN go get github.com/GeertJohan/go.rice/rice
RUN cd server && rice embed-go
RUN echo "Building Version=`git describe --always`"
RUN CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-X main.version=`git describe --always`" -o easybotsvc

FROM alpine
RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Taipei /etc/localtime
RUN echo "Asia/Taipei" > /etc/timezone
RUN apk --no-cache add ca-certificates

WORKDIR /easybot
COPY --from=build-env /go/src/easybot/easybotsvc /easybot/easybotsvc
RUN chmod 777 easybotsvc
CMD ["./easybotsvc"]
