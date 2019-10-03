FROM golang:latest

MAINTAINER luoxiaojun1992 <luoxiaojun1992@sina.cn>

WORKDIR /go/src/app

COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

EXPOSE 8888

CMD ["sh", "-c", "cd", "/go/src/app/src", "&&", "go", "run", "main.go"]
