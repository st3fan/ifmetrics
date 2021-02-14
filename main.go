// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var interfaceNames = []string{
	"vlan1",
	"vlan2",
	"br0",
}

var metricNames = []string{
	"tx_bytes",
	"tx_packets",
	"rx_bytes",
	"rx_packets",
}

type Settings struct {
	Server    string `envconfig:"INFLUXDB_SERVER"`
	Token     string `envconfig:"INFLUXDB_TOKEN"`
	Bucket    string `envconfig:"INFLUXDB_BUCKET"`
	Org       string `envconfig:"INFLUXDB_ORG"`
	Hostname  string `encvonfig:"HOSTNAME"`
	Interface string `envconfig:"INTERFACE"`
	Interval  int    `envconfig:"INTERVAL" default:"15"`
}

func newSettingsFromEnv() (Settings, error) {
	var settings Settings
	if err := envconfig.Process("", &settings); err != nil {
		return Settings{}, errors.Wrap(err, "failed to parse settings from environment")
	}
	return settings, nil
}

type Application struct {
	settings Settings
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func newApplication(settings Settings) (*Application, error) {
	client := influxdb2.NewClient(settings.Server, settings.Token)
	writeAPI := client.WriteAPI(settings.Org, settings.Bucket)

	return &Application{
		settings: settings,
		client:   client,
		writeAPI: writeAPI,
	}, nil
}

func (app *Application) readInterfaceMetric(ifname string, metric string) string {
	body, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/statistics/%s", ifname, metric))
	if err != nil {
		return "0"
	}
	return string(body[0 : len(body)-1])
}

func (app *Application) collectInterfaceMetrics(ifname string, metricNames []string) string {
	metrics := ""
	for _, metricName := range metricNames {
		if len(metrics) != 0 {
			metrics += ","
		}
		metrics += metricName
		metrics += "="
		metrics += app.readInterfaceMetric(ifname, metricName)
	}
	return metrics
}

func (app *Application) recordInterfaceMetrics() {
	now := time.Now().UnixNano()
	for _, interfaceName := range interfaceNames {
		metrics := app.collectInterfaceMetrics(interfaceName, metricNames)
		record := fmt.Sprintf("wrtstats,host=%s,interface=%s %s %d",
			app.settings.Hostname, interfaceName, metrics, now)
		app.writeAPI.WriteRecord(record)
	}
	app.writeAPI.Flush()
}

func (app *Application) Run() {
	ticker := time.NewTicker(time.Duration(app.settings.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			app.recordInterfaceMetrics()
		}
	}
}

func main() {
	settings, err := newSettingsFromEnv()
	if err != nil {
		log.Fatal("[F] Failed to create settings:", err)
	}

	app, err := newApplication(settings)
	if err != nil {
		log.Fatal("[F] Failed to create application:", err)
	}

	app.Run()
}
