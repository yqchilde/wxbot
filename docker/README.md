### Run

```shell
docker run -d \
    --name="wxbot" \
    --restart=always \
    -v $(pwd)/config.yaml:/app/config.yaml \
    yqchilde/wxbot
```