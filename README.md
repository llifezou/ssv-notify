# ssv-notify

Monitor ssv cluster operator.

Monitoring data sources: ssv api and ssvscan api.

## start

### Require

need [Go 1.20.10 or higher](https://golang.org/dl/):

```
wget -c https://golang.org/dl/go1.20.10.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

### Build

```
git clone git@github.com:llifezou/ssv-notify.git
cd ssv-notify
go build -o ssv-notify main.go
```

### Use

- Supports lark and telegram alarms, one must be configured
- Check every 6.4 minutes
- When the API is abnormal and the operator is inactive, an alarm will be sent.
- Support monitoring multiple clusters
- Supports monitoring of all operators or specified operators in the cluster



Edit `config/config.yaml`

```
network: holesky
larkconfig:
  webhook: https://open.larksuite.com/open-apis/bot/v2/hook/e836xxxxx-xxx-xxx-xxxxxx
telegramconfig:
  accesstoken: 665xxxxx:AAAXXXXXXXX-XXX
  chatid: -413xxxxxxx
aim: "all" # 23,25
clusterowner:
  - 0x34415xxxx
  - 0xdBCC5xxxx
```



Run

```
nohup ./ssv-notify run --config ./config/config.yaml &
```

