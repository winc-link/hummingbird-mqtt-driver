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
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cast"
	"github.com/winc-link/hummingbird-mqtt-driver/config"
	constants "github.com/winc-link/hummingbird-mqtt-driver/constant"
	"github.com/winc-link/hummingbird-mqtt-driver/dtos"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/deviceonline"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/server"
	"github.com/winc-link/hummingbird-mqtt-driver/mqttclient"
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
	dr.sd.GetLogger().Infof("device notify %s %v", deviceId, device)
	switch t {
	case commons.DeviceAddNotify:
		deviceonline.AddDeviceStatusTimeController(device)
	case commons.DeviceDeleteNotify:
		deviceonline.DelDeviceStatusTimeController(deviceId)
	}
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr MQTTProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	dr.sd.GetLogger().Infof("product notify %s %v", productId, product)
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
	propertySet.MsgId = data.MsgId
	propertySet.Version = data.Version
	propertySet.Data = data.Data
	var topic string
	if product.NodeType == commons.NodeTypeGateway || product.NodeType == commons.NodeTypeDevice {
		topic = fmt.Sprintf(constants.TopicDevicePropertySet, deviceId)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDevicePropertySet, deviceId)
	}
	dr.mqttClient.Publish(topic, 1, false, propertySet.Marshal())
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
		topic = fmt.Sprintf(constants.TopicDevicePropertyQuery, deviceId)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDevicePropertyQuery, deviceId)
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
		topic = fmt.Sprintf(constants.TopicDeviceServiceInvoke, deviceId)
	} else if product.NodeType == commons.NodeTypeSubDevice {
		topic = fmt.Sprintf(constants.TopicSubDeviceServiceInvoke, deviceId)
	}
	dr.mqttClient.Publish(topic, 1, false, propertySet.Marshal())
	return nil
}

func (dr MQTTProtocolDriver) HandlePropertyReportDebug(ctx context.Context, deviceId string, data model.PropertyReport) error {
	data.Time = time.Now().UnixMilli()
	newData := make(map[string]interface{})
	for k, v := range data.Data {
		newData[k] = cast.ToFloat64(v)
	}
	data.Data = newData
	resp, _ := dr.sd.PropertyReport(deviceId, data)
	if resp.Success != true {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

func (dr MQTTProtocolDriver) HandleEventReportDebug(ctx context.Context, deviceId string, data model.EventReport) error {
	data.Time = time.Now().UnixMilli()
	resp, _ := dr.sd.EventReport(deviceId, data)
	if resp.Success != true {
		return errors.New(resp.ErrorMessage)
	}
	return nil
}

func (dr MQTTProtocolDriver) GatewayControlSet(ctx context.Context, deviceId string, data model.GatewayControlSet) error {
	token := dr.mqttClient.Publish(fmt.Sprintf(constants.TopicGatewayControlSetTopic, deviceId), 1, false, data)
	if token.Wait() && token.Error() != nil {
		dr.sd.GetLogger().Errorf("gateway control set error: %s", token.Error())
	}
	return nil
}

// NewMQTTProtocolDriver MQTT协议驱动
func NewMQTTProtocolDriver(sd *service.DriverService) *MQTTProtocolDriver {
	cfg := config.GetConfig()
	if cfg.OnlineType == constants.CustomOnline && cfg.Interval > 0 {
		sd.GetLogger().Infof("onlineType: %s,interval: %d", cfg.OnlineType, cfg.Interval)
		deviceonline.InitDeviceStatusTimeController(sd)
		for _, device := range sd.GetDeviceList() {
			deviceonline.AddDeviceStatusTimeController(device)
		}
	}

	go server.NewMQTTService(sd).Start()
	time.Sleep(2 * time.Second)
	return &MQTTProtocolDriver{
		sd:         sd,
		mqttClient: mqttclient.NewMQTTClient(sd),
	}
}
