# Interface Metrics
__Stefan Arentz, February 2021__

This program collects minimal metrics for an interface and then submits it to InfluxDB.


## Compile to run on the Asus RT-AC68

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

> TODO Figure out how to permanently run a program like this on the RT-AC68U.

