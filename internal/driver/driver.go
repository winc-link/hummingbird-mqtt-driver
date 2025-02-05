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
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	constants "github.com/winc-link/hummingbird-mqtt-driver/constant"
	"github.com/winc-link/hummingbird-mqtt-driver/dtos"
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

// DeviceNotify 设备添加/修改/删除通知
func (dr MQTTProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr MQTTProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	return nil
}

// Stop 驱动退出通知。
func (dr MQTTProtocolDriver) Stop(ctx context.Context) error {
	for _, d := range dr.sd.GetDeviceList() {
		err := dr.sd.Offline(d.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// HandlePropertySet 设备属性设置
func (dr MQTTProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.DeviceNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.DeviceNotFound]),
			},
		})
	}
	product, ok := dr.sd.GetProductById(device.ProductId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.ProductNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.ProductNotFound]),
			},
		})
	}
	var propertySet dtos.PropertySet
	propertySet.Id = data.MsgId
	propertySet.Version = data.Version
	propertySet.Params = data.Data
	var topic string
	if product.NodeType == commons.NodeTypeGateway || product.NodeType == commons.NodeTypeDevice {
		topic = fmt.Sprintf(constants.TopicDevicePropertySet, deviceId, product.Id)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDevicePropertySet, deviceId, product.Id)
	}
	dr.mqttClient.Publish(topic, 1, false, propertySet.Marshal())

	_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
		MsgId: data.MsgId,
		Data: model.PropertySetResponseData{
			Success: true,
			Code:    uint32(constants.DefaultSuccessCode),
		},
	})
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr MQTTProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.DeviceNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.DeviceNotFound]),
			},
		})
	}
	product, ok := dr.sd.GetProductById(device.ProductId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.ProductNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.ProductNotFound]),
			},
		})
	}
	var propertySet dtos.PropertyQuery
	propertySet.Id = data.MsgId
	propertySet.Version = data.Version
	propertySet.Params = data.Data
	var topic string
	if product.NodeType == commons.NodeTypeGateway || product.NodeType == commons.NodeTypeDevice {
		topic = fmt.Sprintf(constants.TopicDevicePropertyQuery, deviceId, product.Id)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDevicePropertyQuery, deviceId, product.Id)
	}
	dr.mqttClient.Publish(topic, 1, false, propertySet.Marshal())
	return nil
}

// HandleServiceExecute 设备服务调用
func (dr MQTTProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.DeviceNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.DeviceNotFound]),
			},
		})
	}
	product, ok := dr.sd.GetProductById(device.ProductId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constants.ProductNotFound),
				ErrorMessage: string(constants.ErrorCodeMsgMap[constants.ProductNotFound]),
			},
		})
	}
	var propertySet dtos.ServiceInvoke
	propertySet.Id = data.MsgId
	propertySet.Version = data.Version
	propertySet.Params = data.Data
	var topic string
	if product.NodeType == commons.NodeTypeGateway || product.NodeType == commons.NodeTypeDevice {
		topic = fmt.Sprintf(constants.TopicDeviceServiceInvoke, deviceId, product.Id)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDeviceServiceInvoke, deviceId, product.Id)
	}
	dr.mqttClient.Publish(topic, 1, false, propertySet.Marshal())
	return nil
}

// NewMQTTProtocolDriver MQTT协议驱动
func NewMQTTProtocolDriver(sd *service.DriverService) *MQTTProtocolDriver {
	go server.NewMQTTService(sd).Start()
	time.Sleep(2 * time.Second)
	return &MQTTProtocolDriver{
		sd:         sd,
		mqttClient: client.NewMQTTClient(sd),
	}
}
