#!/bin/sh
remote="user@server"
service="name.service"
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot
ssh $remote "systemctl stop $service"
scp bot $remote:/var/www/itqonbot/bot
ssh $remote "systemctl start $service"
