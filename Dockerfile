FROM golang:1.20.1 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN apt-get update
RUN apt-get install cron -y
RUN apt-get install logrotate -y
RUN apt-get install postgresql postgresql-contrib -y

# Adding crontab to the appropriate location
ADD crontab /etc/cron.d/sms-cron-file

# Giving permission to crontab file
RUN chmod 0644 /etc/cron.d/sms-cron-file

# Running crontab
RUN crontab /etc/cron.d/sms-cron-file

RUN apt-get install -y iputils-ping

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN export PATH=$PATH:$(go env GOPATH)/bin

WORKDIR /opt/sms/

ENTRYPOINT [ "air" ]

