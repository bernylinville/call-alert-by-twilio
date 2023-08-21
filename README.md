# call-alert-by-twilio

```bash
amtool --alertmanager.url=http://127.0.0.1:9093/ alert add alertname="HostOutOfMemory" severity="call" job="test-alert" instance="localhost" exporter="none" cluster="test" hostname="test9" summary="Host out of memory"
```

```yaml
route:
  routes:
    - receiver: call
      match:
        severity: call

receivers:
  - name: 'call'
    webhook_configs:
      - url: 'http://twilio:1337/call'
        send_resolved: false
```