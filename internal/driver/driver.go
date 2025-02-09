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

package driver

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/client"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/server"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"time"
)

type MQTTProtocolDriver struct {
	sd         *service.DriverService
	mqttClient mqtt.Client
}

// CloudPluginNotify 云插件启动/停止通知
func (dr MQTTProtocolDriver) CloudPluginNotify(ctx context.Context, t commons.CloudPluginNotifyType, name string) error {
	//TODO implement me
	panic("implement me")
}

// DeviceNotify 设备添加/修改/删除通知
func (dr MQTTProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	//TODO implement me
	panic("implement me")
}

// ProductNotify 产品添加/修改/删除通知
func (dr MQTTProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	//TODO implement me
	panic("implement me")
}

// Stop 驱动退出通知。
func (dr MQTTProtocolDriver) Stop(ctx context.Context) error {
	for _, device := range dr.sd.GetDeviceList() {
		dr.sd.Offline(device.Id)
	}
	return nil
}

// HandlePropertySet 设备属性设置
func (dr MQTTProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr MQTTProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	//TODO implement me
	panic("implement me")
}

// HandleServiceExecute 设备服务调用
func (dr MQTTProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	//TODO implement me
	panic("implement me")
}

// NewMQTTProtocolDriver MQTT协议驱动
func NewMQTTProtocolDriver(sd *service.DriverService) *MQTTProtocolDriver {
	go server.NewMQTTService(sd).Start()
	time.Sleep(1 * time.Second)
	return &MQTTProtocolDriver{
		sd:         sd,
		mqttClient: client.NewMQTTClient(sd),
	}
}
