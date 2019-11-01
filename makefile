RELEASE=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APP?=app
PORT?=11009
ImageName?=sherrymail
ContainerName?=mail
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

cleanDocker:
	sh clean.sh

clean:
	rm -f ${APP}

build:
	GOOS=linux GOARCH=amd64 go build -tags netgo \
	-ldflags "-s -w -X version.Release=${RELEASE} \
	-X version.Commit=${COMMIT} \
	-X version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

docker: build
	docker build -t ${ImageName} .
	rm -f ${APP}
	docker images

run: docker cleanDocker
	docker run -d --name ${ContainerName} \
	-v /etc/localtime:/etc/localtime:ro \
	-v ${CURDIR}/template:/app/template  \
	-p ${PORT}:80 \
	--env-file ${CURDIR}envfile \
	${ImageName}

stop:
	docker stop ${ContainerName}

log:
	 docker logs -f -t --tail 20 ${ContainerName}

rm:
	docker rm ${ContainerName}
	docker ps -a

login:
	docker exec -it ${ContainerName} /bin/bash

test:
	PORT=11012 MailAccount=andyliu MailPassword=2iduudgR@2019 MailServer=smtp.sinica.edu.tw MailServerPort=25 ./sherrymail
