#!/bin/bash

echo "1. 获取配置文件到config.yaml"
wget -q -O config.yaml https://raw.fastgit.org/yqchilde/wxbot/hook/config.yaml
chmod 755 config.yaml

echo "2. 检查wxbot容器"
docker pull yqchilde/wxbot:latest
docker ps | grep wxbot | awk '{print $1}' | xargs docker rm -f

echo "3. 启动wxbot容器"
docker run -dit --name wxbot -p 9528:9528 -v $(pwd)/config.yaml:/app/config.yaml -v $(pwd)/data:/app/data yqchilde/wxbot:latest
echo "3. 启动完成，日志查看: docker logs -f wxbot"
