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
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/winc-link/hummingbird-mqtt-driver/config"
	constants "github.com/winc-link/hummingbird-mqtt-driver/constant"
	"github.com/winc-link/hummingbird-mqtt-driver/dtos"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/pkg/convert"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"reflect"
	"strconv"
	"time"
)

var driverService *service.DriverService

func NewMQTTClient(sd *service.DriverService) mqtt.Client {
	driverService = sd
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", "0.0.0.0", 1883))
	opts.SetClientID(constants.MQTTInnerClientId)
	opts.SetUsername(constants.MQTTInnerUsername)
	opts.SetPassword(constants.MQTTInnerPassword)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	client.Subscribe(constants.TopicDevicePropertyReport, 1, propertyPostCallback)
	client.Subscribe(constants.TopicDeviceEventReport, 1, eventPostCallback)
	client.Subscribe(constants.TopicDevicePropertySetReply, 1, propertySetReplyCallback)
	client.Subscribe(constants.TopicDevicePropertyQueryReply, 1, propertyQueryReplyCallback)
	client.Subscribe(constants.TopicDeviceServiceInvokeReply, 1, serviceInvokeReplyCallback)
	return client
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	driverService.GetLogger().Infof("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	driverService.GetLogger().Infof("Connect lost: %v", err)
}

// serviceInvokeCallback 设备服务调用响应
func serviceInvokeReplyCallback(client mqtt.Client, message mqtt.Message) {
	topic := dtos.Topic(message.Topic())
	deviceId := topic.GetThingModelTopicDeviceId()
	var serviceInvokeReply dtos.ServiceInvokeReply
	err := json.Unmarshal(message.Payload(), &serviceInvokeReply)
	if err != nil {
		driverService.GetLogger().Errorf("the format of result is error %s", err.Error())
		return
	}
	if err = driverService.ServiceExecuteResponse(deviceId, model.ServiceExecuteResponse{
		MsgId: serviceInvokeReply.Id,
		Data:  serviceInvokeReply.Params,
	}); err != nil {
		driverService.GetLogger().Errorf("error %s", err.Error())
	}
}

// propertyQueryCallback
func propertyQueryReplyCallback(client mqtt.Client, message mqtt.Message) {
	topic := dtos.Topic(message.Topic())
	deviceId := topic.GetThingModelTopicDeviceId()
	var propertyQueryReply dtos.PropertyQueryReply
	err := json.Unmarshal(message.Payload(), &propertyQueryReply)
	if err != nil {
		driverService.GetLogger().Errorf("the format of result is error %s", err.Error())
		return
	}
	if err = driverService.PropertyGetResponse(deviceId, model.PropertyGetResponse{
		MsgId: propertyQueryReply.Id,
		Data:  propertyQueryReply.Params,
	}); err != nil {
		driverService.GetLogger().Errorf("error %s", err.Error())
	}
}

// propertySetReplyCallback 设置设备属性响应
func propertySetReplyCallback(client mqtt.Client, message mqtt.Message) {
	topic := dtos.Topic(message.Topic())
	deviceId := topic.GetThingModelTopicDeviceId()
	var propertySetReply dtos.PropertySetReply
	err := json.Unmarshal(message.Payload(), &propertySetReply)
	if err != nil {
		driverService.GetLogger().Errorf("the format of result is error %s", err.Error())
		return
	}
	if err = driverService.PropertySetResponse(deviceId, model.PropertySetResponse{
		MsgId: propertySetReply.Id,
		Data:  propertySetReply.Params,
	}); err != nil {
		driverService.GetLogger().Errorf("error %s", err.Error())
	}
}

// eventPostCallback 设备事件上报
func eventPostCallback(client mqtt.Client, message mqtt.Message) {
	topic := dtos.Topic(message.Topic())
	deviceId := topic.GetThingModelTopicDeviceId()
	productId := topic.GetThingModelTopicProductId()
	var eventPost dtos.EventPost
	err := json.Unmarshal(message.Payload(), &eventPost)

	var eventPostReply dtos.CommonResponse
	eventPostReply.Id = eventPost.Id
	eventPostReply.Version = eventPost.Version

	if err != nil {
		eventPostReply.Code = int(constants.FormatErrorCode)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.FormatErrorCode])
		client.Publish(fmt.Sprintf(constants.TopicDeviceEventReportReply, deviceId, productId), 1, false, eventPostReply.Marshal())
		return
	}

	// 检查设备是否存在
	if _, ok := driverService.GetDeviceById(deviceId); !ok {
		eventPostReply.Code = int(constants.DeviceNotFound)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DeviceNotFound])
		client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, eventPostReply.Marshal())
		return
	}

	// 检查产品是否存在
	if _, ok := driverService.GetProductById(productId); !ok {
		eventPostReply.Code = int(constants.ProductNotFound)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.ProductNotFound])
		client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, eventPostReply.Marshal())
		return
	}

	// 事件code检测
	_, ok := driverService.GetProductEventByCode(productId, eventPost.Params.EventCode)
	if !ok {
		eventPostReply.Code = int(constants.EventCodeNotFound)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.EventCodeNotFound]) + fmt.Sprintf(" %s is undefined", eventPost.Params.EventCode)
		client.Publish(fmt.Sprintf(constants.TopicDeviceEventReportReply, deviceId, productId), 1, false, eventPostReply.Marshal())
		return
	}

	if eventPost.Params.EventTime == 0 {
		eventPost.Params.EventTime = time.Now().UnixMilli()
	}
	_, err = driverService.EventReport(deviceId, model.NewEventReport(eventPost.Sys.Ack, eventPost.Params))
	if !eventPost.Sys.Ack {
		return
	}
	if err != nil {
		eventPostReply.Code = int(constants.SystemErrorCode)
		eventPostReply.Success = false
		eventPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + " :" + err.Error()
	} else {
		eventPostReply.Code = int(constants.DefaultSuccessCode)
		eventPostReply.Success = true
		eventPostReply.ErrorMessage = ""
	}
	client.Publish(fmt.Sprintf(constants.TopicDeviceEventReportReply, deviceId, productId), 1, false, eventPostReply.Marshal())
}

// propertyPostCallback 设备属性上报
func propertyPostCallback(client mqtt.Client, message mqtt.Message) {
	topic := dtos.Topic(message.Topic())
	deviceId := topic.GetThingModelTopicDeviceId()
	productId := topic.GetThingModelTopicProductId()
	var propertyPost dtos.PropertyPost
	err := json.Unmarshal(message.Payload(), &propertyPost)
	var propertyPostReply dtos.CommonResponse
	propertyPostReply.Id = propertyPost.Id
	propertyPostReply.Version = propertyPost.Version

	if err != nil {
		propertyPostReply.Code = int(constants.FormatErrorCode)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.FormatErrorCode])
		client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
		return
	}

	// 检查设备是否存在
	if _, ok := driverService.GetDeviceById(deviceId); !ok {
		propertyPostReply.Code = int(constants.DeviceNotFound)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.DeviceNotFound])
		client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
		return
	}

	// 检查产品是否存在
	if _, ok := driverService.GetProductById(productId); !ok {
		propertyPostReply.Code = int(constants.ProductNotFound)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.ProductNotFound])
		client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
		return
	}

	var delPropertyCode []string
	// code 、time、max～min、len检测
	for code, param := range propertyPost.Params {
		if property, ok := driverService.GetProductPropertyByCode(productId, code); !ok {
			//推送一条错误消息到客户端
			delPropertyCode = append(delPropertyCode, code)
			propertyPostReply.Code = int(constants.PropertyCodeNotFound)
			propertyPostReply.Success = false
			propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.PropertyCodeNotFound]) + fmt.Sprintf(" %s is undefined", code)
			client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
			continue
		} else {
			value := param.Value
			if config.GetConfig().TslParamVerify {
				//推送一条错误消息到客户端
				if verifyErrorCode, verifyErrorMsg := verifyParam(property, param); verifyErrorCode != constants.DefaultSuccessCode {
					delPropertyCode = append(delPropertyCode, code)
					propertyPostReply.Code = int(verifyErrorCode)
					propertyPostReply.Success = false
					propertyPostReply.ErrorMessage = string(verifyErrorMsg)
					client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
					continue
				}
			}
			if param.Time == 0 {
				propertyPost.Params[code] = model.PropertyData{
					Time:  time.Now().UnixMilli(),
					Value: value,
				}
			}
		}
	}
	filterPropertyPost := propertyPost.Params
	for _, code := range delPropertyCode {
		delete(filterPropertyPost, code)
	}
	propertyPost.Params = filterPropertyPost
	if len(propertyPost.Params) == 0 {
		return
	}
	_, err = driverService.PropertyReport(deviceId, model.NewPropertyReport(propertyPost.Sys.Ack, propertyPost.Params))
	if !propertyPost.Sys.Ack {
		return
	}
	if err != nil {
		propertyPostReply.Code = int(constants.SystemErrorCode)
		propertyPostReply.Success = false
		propertyPostReply.ErrorMessage = string(constants.ErrorCodeMsgMap[constants.SystemErrorCode]) + " :" + err.Error()
	} else {
		propertyPostReply.Code = int(constants.DefaultSuccessCode)
		propertyPostReply.Success = true
		propertyPostReply.ErrorMessage = ""
	}
	client.Publish(fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId), 1, false, propertyPostReply.Marshal())
}

func verifyParam(property model.Property, param model.PropertyData) (constants.ErrorCode, constants.ErrorMessage) {
	switch property.TypeSpec.Type {
	case "int":
		var intOrFloatSpecs dtos.IntOrFloatSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &intOrFloatSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToInt(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		min, _ := convert.StringToInt(intOrFloatSpecs.Min)
		max, _ := convert.StringToInt(intOrFloatSpecs.Max)

		if t < min || t > max {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "float":
		var intOrFloatSpecs dtos.IntOrFloatSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &intOrFloatSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an float type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		min, _ := convert.StringToFloat64(intOrFloatSpecs.Min)
		max, _ := convert.StringToFloat64(intOrFloatSpecs.Max)
		if t < min || t > max {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "bool":
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		if !(t == 0 || t == 1) {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	case "text":
		var textSpecs dtos.TextSpecs
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &textSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToString(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an string type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		length, err := strconv.Atoi(textSpecs.Length)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		if length < len(t) {
			return constants.ReportDataLengthErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataLengthErrorCode]
		}
	case "enum":
		enumSpecs := make(map[string]string)
		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &enumSpecs)
		if err != nil {
			return constants.SystemErrorCode,
				constants.ErrorCodeMsgMap[constants.SystemErrorCode]
		}
		t, err := convert.GetInterfaceToFloat64(param.Value)
		if err != nil {
			return constants.PropertyReportTypeErrorCode,
				constants.ErrorCodeMsgMap[constants.PropertyReportTypeErrorCode] +
					constants.ErrorMessage(fmt.Sprintf(": %s value is a %s type and not an int type", property.Code, reflect.TypeOf(param.Value).Kind()))
		}
		if _, ok := enumSpecs[strconv.Itoa(int(t))]; !ok {
			return constants.ReportDataRangeErrorCode,
				constants.ErrorCodeMsgMap[constants.ReportDataRangeErrorCode]
		}
	}
	return constants.DefaultSuccessCode,
		constants.ErrorCodeMsgMap[constants.DefaultSuccessCode]
}
