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

package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	constants "github.com/winc-link/hummingbird-mqtt-driver/constant"
	"github.com/winc-link/hummingbird-mqtt-driver/dtos"
	"github.com/winc-link/hummingbird-mqtt-driver/internal/pkg/tool"
	"net"
)

// OnAccept TCP连接建立时调用
// TCP连接限速，黑白名单等.
func OnAccept(ctx context.Context, conn net.Conn) bool {
	return true
}

// OnStop 当gmqtt退出时调用
func OnStop(ctx context.Context) {}

// OnSubscribe 收到订阅请求时调用
// 校验订阅是否合法
func OnSubscribe(ctx context.Context, client server.Client, req *server.SubscribeRequest) error {
	if client.ClientOptions().ClientID == constants.MQTTInnerClientId {
		return nil
	}
	for _, topic := range req.Subscribe.Topics {
		t := dtos.Topic(string(topic.Name))
		deviceId := t.GetThingModelTopicDeviceId()
		productId := t.GetThingModelTopicProductId()
		if deviceId == "" || productId == "" {
			return errors.New("subscribe Unauthorized")
		}
		dev, ok := GlobalDriverService.GetDeviceById(deviceId)
		if !ok {
			return errors.New("subscribe Unauthorized")
		}
		product, ok := GlobalDriverService.GetProductById(productId)
		if !ok {
			return errors.New("subscribe Unauthorized")
		}

		if dev.ProductId != product.Id {
			return errors.New("subscribe Unauthorized")
		}
	}
	// 如果用户想检查订阅消息时topic是否合法，请用户自行完成下面相关内容
	//	if topic.Name == fmt.Sprintf(constants.TopicDevicePropertyReportReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicDeviceEventReportReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicDevicePropertySet, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicDevicePropertyQuery, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicDeviceServiceInvoke, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDeviceOnlineReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDeviceOfflineReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDevicePropertyReportReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDeviceEventReportReply, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDevicePropertySet, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDevicePropertyQuery, deviceId, productId) ||
	//		topic.Name == fmt.Sprintf(constants.TopicSubDeviceServiceInvoke, deviceId, productId) {
	//		continue
	//	}
	//	return errors.New("subscribe Unauthorized")
	//}

	return nil
}

// OnSubscribed 订阅成功后调用
// 统计订阅报文数量
func OnSubscribed(ctx context.Context, client server.Client, subscription *gmqtt.Subscription) {}

// OnUnsubscribe 取消订阅时调用
// 校验是否允许取消订阅
func OnUnsubscribe(ctx context.Context, client server.Client, req *server.UnsubscribeRequest) error {
	return nil
}

// OnUnsubscribed 取消订阅成功后调用
// 统计订阅报文数
func OnUnsubscribed(ctx context.Context, client server.Client, topicName string) {}

// OnMsgArrived 收到消息发布报文时调用
// 校验发布权限，改写发布消息
func OnMsgArrived(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {
	if client.ClientOptions().ClientID == constants.MQTTInnerClientId {
		return nil
	}
	topic := dtos.Topic(string(req.Publish.TopicName))
	deviceId := topic.GetThingModelTopicDeviceId()
	productId := topic.GetThingModelTopicProductId()
	device, ok := GlobalDriverService.GetDeviceById(deviceId)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	product, ok := GlobalDriverService.GetProductById(productId)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if device.ProductId != product.Id {
		return fmt.Errorf("unauthorized")
	}
	// 如果用户想检查收到消息时topic是否合法，请用户自行完成下面相关内容
	//if string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/property/post", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/event/post", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/property/query_reply", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/property/set_reply", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/service/invoke_reply", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/online", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/offline", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/property/post", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/event/post", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/property/set_reply", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/property/query_reply", deviceId, productId) ||
	//	string(topic) == fmt.Sprintf(constants.TopicPrefix+"%s/%s/thing/sub/service/invoke_reply", deviceId, productId) {
	//	return nil
	//}
	// todo 消息改写。可以在物联网页面使用javascript或者php书写一段消息解析代码，然后把这段代码下发到驱动中，驱动根据这段代码动态解析设备上报数据，可以使用如下仓库来实现此功能。
	//java script
	//https://github.com/robertkrimen/otto
	//php
	//https://github.com/deuill/go-php
	//go ?
	return nil
}

// OnBasicAuth 收到连接请求报文时调用
// 客户端连接鉴权
func OnBasicAuth(ctx context.Context, client server.Client, req *server.ConnectRequest) (err error) {
	clientId := string(req.Connect.ClientID)
	username := string(req.Connect.Username)
	password := string(req.Connect.Password)
	if clientId == "" || username == "" || password == "" {
		return fmt.Errorf("unauthorized")
	}
	if clientId == constants.MQTTInnerClientId && username == constants.MQTTInnerUsername && password == constants.MQTTInnerPassword {
		return nil
	}
	dev, ok := GlobalDriverService.GetDeviceById(clientId)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	product, ok := GlobalDriverService.GetProductById(dev.ProductId)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if username != (dev.Id + "&" + product.Key) {
		return fmt.Errorf("unauthorized")
	}
	if password != tool.HmacMd5(dev.Secret, dev.Id+"&"+product.Key) {
		return fmt.Errorf("unauthorized")
	}
	err = GlobalDriverService.Online(clientId)
	if err != nil {
		GlobalDriverService.GetLogger().Errorf("device online err:%s", err.Error())
		return err
	}

	return nil
}

// OnEnhancedAuth 收到带有AuthMeth的连接请求报文时调用（V5特性）
// 客户端连接鉴权
func OnEnhancedAuth(ctx context.Context, client server.Client, req *server.ConnectRequest) (resp *server.EnhancedAuthResponse, err error) {
	return
}

// OnReAuth 收到Auth报文时调用（V5特性）
// 客户端连接鉴权
func OnReAuth(ctx context.Context, client server.Client, auth *packets.Auth) (*server.AuthResponse, error) {
	return nil, nil
}

// OnConnected 客户端连接成功后调用
// 统计在线客户端数量
func OnConnected(ctx context.Context, client server.Client) {}

// OnSessionCreated
// 统计session数量
func OnSessionCreated(ctx context.Context, client server.Client) {}

// OnSessionResumed 客户端从旧session恢复后调用
// 统计session数量
func OnSessionResumed(ctx context.Context, client server.Client) {}

// OnSessionTerminated session删除后调用
// 统计session数量
func OnSessionTerminated(ctx context.Context, clientID string, reason server.SessionTerminatedReason) {
}

// OnDelivered 消息从broker投递到客户端后调用
func OnDelivered(ctx context.Context, client server.Client, msg *gmqtt.Message) {}

// OnClosed 统计在线客户端数量
func OnClosed(ctx context.Context, client server.Client, err error) {
	clientId := client.ClientOptions().ClientID
	if clientId != constants.MQTTInnerClientId {
		err = GlobalDriverService.Offline(clientId)
		if err != nil {
			GlobalDriverService.GetLogger().Errorf("device offline err:%s", err.Error())
		}
	}

}

// OnMsgDropped 消息被丢弃时调用
func OnMsgDropped(ctx context.Context, clientID string, msg *gmqtt.Message, err error) {}

// OnWillPublish 发布遗嘱消息前
// 修改或丢弃遗嘱消息
func OnWillPublish(ctx context.Context, clientID string, req *server.WillMsgRequest) {}

// OnWillPublished 发布遗嘱消息后
func OnWillPublished(ctx context.Context, clientID string, msg *gmqtt.Message) {}
