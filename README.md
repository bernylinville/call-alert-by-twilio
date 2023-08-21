# call-alert-by-twilio

```yaml
route:
  routes:
    - receiver: call
      match:
        severity: call

receivers:
  - name: 'call'
    webhook_configs:
      - url: 'http://twilio:1337/answer'
        send_resolved: false
```
