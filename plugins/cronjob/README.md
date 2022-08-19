# Cronjob

定时任务

## 配置

```yaml
cronjob:
  enable: true
  myb:
    cron: "0 30 9 * * *"
    groups: [ "Test" ]
```