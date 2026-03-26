/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package main

import (
	"github.com/winc-link/hummingbird-mqtt-driver/config"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/driver"
	"github.com/winc-link/hummingbird-sdk-go/constants"
	"github.com/winc-link/hummingbird-sdk-go/datadb/influxdb"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

func main() {
	//driverService := service.NewDriverService("official-mqtt-driver-v3.1")

	driverService := service.NewDriverService("test",
		service.WithCustomRedisBasesConfig(&service.RedisBasesConnConfig{
			Address:  "*.*.:6379",
			Password: "",
			DB:       0,
		}),
		service.WithCustomDataBasesConfig(&service.DataBasesConnConfig{
			Type: constants.DataBasesInfluxdb,
			InfluxDB: influxdb.DbClient{
				Org:       "hummingbird",
				Bucket:    "device-data",
				LogBucket: "device-log",
				Url:       "http://*.221.36.14:8086",
				Token:     "6S4LFh_kP0-RHqIYjYlXgvGfXOgMIkqMkinZDePKiXbcmgIQzQcm5mV5GfSQEDVqcKzQ5WIixO7AEmwHQ17JmQ==",
			},

			//Type: constants.DataBasesTdengine,
			//Tdengine: tdengine.DbClient{
			//	Dsn: "root:taosdata@ws(127.0.0.1:6041)/devicedata",
			//},
		}), service.WithCustomMetaBasesConfig(&service.MetaBasesConnConfig{
			Type: constants.MetadataMysql,
			Dns:  "root:!@#12345678.@tcp(*.221.36.14:3306)/hummingbird?charset=utf8mb4&parseTime=True&loc=Local&timeout=2s",
		}))
	config.InitConfig(driverService)
	mqttDriver := driver.NewMQTTProtocolDriver(driverService)
	if err := driverService.Start(mqttDriver); err != nil {
		driverService.GetLogger().Error("driver service start error: %s", err)
		return
	}
}
