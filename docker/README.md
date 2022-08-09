### Run

```shell
docker run -d \
    --name="wxbot" \
    --restart=always \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v $(pwd)/holiday.json:/app/holiday.json \
    -v $(pwd)/imgs:/app/imgs \
    yqchilde/wxbot
```