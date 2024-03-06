# ssv-notify

- Monitor the status of the SSV cluster operator and alarm when the operator is inactive
  - WARN: Monitoring accuracy relies heavily on third-party API. (SSV API & SSVSCAN API)
- Monitor the SSV cluster balance, alarm before liquidation
  - INFO: Does not rely on third-party API, Data comes from eth full node.
- Scan all clusters
  - INFO: Does not rely on third-party API, Data comes from eth full node.

## Start

**Require**

[Go 1.20.10 or higher](https://golang.org/dl/)

```
wget -c https://golang.org/dl/go1.20.10.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

**Build**

```
make build
```

**Function**

- Support mainnet, goerli, holesky
- Supports lark, telegram, gmail, and discord alarms. You must fill in one. If you fill in multiple, multiple alarms will be sent. create an alarm robot:
  - [Create a telegram bot](https://telegram.me/botfather)
  - [Create lark webhook](https://open.larksuite.com/document/client-docs/bot-v3/add-custom-bot)
  - [Create discord webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
  - [Gmail app password](https://support.google.com/accounts/answer/185833)
- Operator monitor
  - Operator status monitoring is checked every 6.4 minutes
  - Support monitoring multiple clusters
  - Supports monitoring of all operators or specified operators in the cluster
- Liquidation monitor
  - Liquidation monitoring is checked every 6 hours
  - Support monitoring multiple clusters
- Cluster scan
  - Scan all cluster data and calculate liquidation date
- Alarm test
  - Send an alarm test to confirm that the alarm is available

**Command**

```
ssv monitoring notifications.

Usage:
  ssv-notify [command]

Available Commands:
  alarm-test          Test alarms can be used
  completion          Generate the autocompletion script for the specified shell
  help                Help about any command
  liquidation-monitor liquidation monitor
  operator-monitor    operator monitor
  ssv-tools           ssv tools

Flags:
  -c, --config string   Path to configuration file (default "./config/config.yaml")
  -h, --help            help for ssv-notify

Use "ssv-notify [command] --help" for more information about a command.
```

**config.yaml**

```
# support: mainnet / goerli / holesky
network: mainnet

# eth execution layer rpc endpoint
ethrpc: # https://mainnet.infura.io/v3/xxxxxx

# If filled in, the alarm will be sent to lark
larkconfig:
  webhook: # https://open.larksuite.com/open-apis/bot/v2/hook/e836xxxxx-xxx-xxx-xxxxxx

# If filled in, the alarm will be sent to telegram
telegramconfig:
  accesstoken: # 665xxxxx:AAAXXXXXXXX-XXX
  chatid: # -413xxxxxxx

# If filled in, the alarm will be sent to the mailbox
gmailconfig:
  from: # xxx@gmail.com
  password: # 'ydtd xxxx xxxx tdxj'
  to: # xxx@xxxmail.com

# If filled in, the alarm will be sent to discord
discordconfig:
  webhook: # https://discord.com/api/webhooks/1214540203973414932/lC9-Wxp3BQ_WxlbOBLCXxxxxxXXXX

# operator monitor: Monitor the operator status of the validator cluster
operatormonitor:
  aim: "all" # Use 'all' to monitor all operators; use commas to monitor some operators, such as 23,25
  clusterowner:
    - # 0x344152eD7110694B004962CD61ddA876559Fd8a4
    - # 0xdBCC5c776E4Ca9AdFBE9f2C341bB32e05f582448

# liquidation monitor: Monitor the balance of the cluster
liquidationmonitor:
  threshold: 30 # Operational Runway less than 30 days will trigger an alarm
  clusterowner:
    - # 0x344152eD7110694B004962CD61ddA876559Fd8a4
    - # 0xdBCC5c776E4Ca9AdFBE9f2C341bB32e05f582448
```

**Run**

After config.yaml is configured correctly

- Operator monitor

  ```
  nohup ./ssv-notify operator-monitor -c ./config/config.yaml > ./operator-monitor.log 2>&1 &
  ```

- Liquidation monitor

  ```
  nohup ./ssv-notify liquidation-monitor -c ./config/config.yaml > ./liquidation-monitor.log 2>&1
  ```

- Cluster scan

  ```
  ./ssv-notify ssv-tools scan-cluster -d ./ -c ./config/config.yaml
  ```

- Alarm test

  ```
  ./ssv-notify alarm-test
  ```

  

### TODO

- Remove third-party API dependencies
