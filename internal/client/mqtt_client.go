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

package client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/winc-link/hummingbird-mqtt-driver/constants"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

var driverService *service.DriverService

func NewMQTTClient(sd *service.DriverService) mqtt.Client {
	driverService = sd
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", "0.0.0.0", 1883))
	opts.SetClientID(constants.MQTTInnerClientId)
	opts.SetUsername(constants.MQTTInnerUsername)
	opts.SetPassword(constants.MQTTInnerPassword)
	opts.SetAutoReconnect(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//todo your MQTT subscription logic code

	return client
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	driverService.GetLogger().Infof("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	driverService.GetLogger().Infof("Connect lost: %v", err)
}
