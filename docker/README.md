### Run

```shell
docker run -d \
    --name="wxbot" \
    --restart=always \
    -p 9528:9528 \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v $(pwd)/data:/app/data \
    yqchilde/wxbot
```
