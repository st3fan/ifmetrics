# Interface Metrics
_Stefan Arentz, February 2021_

This program collects minimal metrics for an interface and then submits it to InfluxDB. This is something that `Telegraf` does pretty well, but I've found that too heavyweight to run on small devices like Wi-Fi routers. This program does the minimal thing needed and then submits to _InfluxDB_ using the simple line format.

I'm running this on an Asus RT-AC68U router through the [AsusWRT-Merlin-NG](https://github.com/RMerl/asuswrt-merlin.ng) project. Some notes below on how to run it more permanently on your router.

![Alt text](screenshot.png?raw=true "Screenshot")


## Compile to run on the Asus RT-AC68U

```
GOOS=linux GOARCH=arm GOARM=5 go build
```

## Configure and run

```
export INFLUXDB_SERVER="https://us-east-1-1.aws.cloud2.influxdata.com"
export INFLUXDB_TOKEN="... YOUR TOKEN ..."
export INFLUXDB_BUCKET="... YOUR BUCKET NAME ..."
export INFLUXDB_ORG="... YOUR ORG NAME ..."

export HOSTNAME="192.168.0.1"
export INTERFACE="vlan2"
export INTERVAL=15

./ifmetrics
```

## Running it as a daemon

I am using Merlin on my Asus RT-AC68U which lets me run this in a `tmux` session. I'd love to turn this into a `.opkg` though, so that you can install and run it more easily as a daemon.

