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

package constants

const (
	MQTTInnerClientId = "hummingbird-mqtt-driver"
	MQTTInnerUsername = "hummingbird-mqtt-driver"
	MQTTInnerPassword = "4475f3a4-ff7d-4d0e-a614-572a8774421b"

	TopicPrefix = "/sys/"

	// 属性上报
	TopicDevicePropertyReport      = TopicPrefix + "+/+/thing/property/post"         //设备->平台 设备属性上报
	TopicDevicePropertyReportReply = TopicPrefix + "%s/%s/thing/property/post_reply" //平台->设备 属性上报响应

	// 事件上报
	TopicDeviceEventReport      = TopicPrefix + "+/+/thing/event/post"         //设备->平台 设备事件上报
	TopicDeviceEventReportReply = TopicPrefix + "%s/%s/thing/event/post_reply" //平台->设备 事件上报响应

	// 设置属性
	TopicDevicePropertySetReply = TopicPrefix + "+/+/thing/property/set_reply" //设备->平台 设置设备属性响应
	TopicDevicePropertySet      = TopicPrefix + "%s/%s/thing/property/set"     //平台->设备 设置设备属性

	// 设备属性查询
	TopicDevicePropertyQueryReply = TopicPrefix + "+/+/thing/property/query_reply" //设备->平台 设备属性查询响应
	TopicDevicePropertyQuery      = TopicPrefix + "%s/%s/thing/property/query"     //平台->设备 设备属性查询

	// 设备服务调用
	TopicDeviceServiceInvokeReply = TopicPrefix + "+/+/thing/service/invoke_reply" //设备->平台 设备服务调用响应
	TopicDeviceServiceInvoke      = TopicPrefix + "%s/%s/thing/service/invoke"     //平台->设备 设备服务调用

	// 子设备上线
	TopicSubDeviceOnline      = TopicPrefix + "+/+/thing/sub/online"         //设备->平台 子设备上线
	TopicSubDeviceOnlineReply = TopicPrefix + "%s/%s/thing/sub/online_reply" //平台->设备 子设备上线响应

	// 子设备下线
	TopicSubDeviceOffline      = TopicPrefix + "+/+/thing/sub/offline"         //设备->平台 子设备下线
	TopicSubDeviceOfflineReply = TopicPrefix + "%s/%s/thing/sub/offline_reply" //平台->设备 子设备下线响应

	//子设备属性上报
	TopicSubDevicePropertyReport      = TopicPrefix + "+/+/thing/sub/property/post"         //设备->平台 子设备属性上报
	TopicSubDevicePropertyReportReply = TopicPrefix + "%s/%s/thing/sub/property/post_reply" //平台->设备 子设备属性上报响应

	//子设备事件上报
	TopicSubDeviceEventReport      = TopicPrefix + "+/+/thing/sub/event/post"         //设备->平台 设备事件上报
	TopicSubDeviceEventReportReply = TopicPrefix + "%s/%s/thing/sub/event/post_reply" //平台->设备 事件上报响应

	//子设备设置属性
	TopicSubDevicePropertySetReply = TopicPrefix + "+/+/thing/sub/property/set_reply" //设备->平台 设置设备属性响应
	TopicSubDevicePropertySet      = TopicPrefix + "%s/%s/thing/sub/property/set"     //平台->设备 设置设备属性

	//子设备设备属性查询
	TopicSubDevicePropertyQueryReply = TopicPrefix + "+/+/thing/sub/property/query_reply" //设备->平台 设备属性查询响应
	TopicSubDevicePropertyQuery      = TopicPrefix + "%s/%s/thing/sub/property/query"     //平台->设备 设备属性查询

	//子设备设备服务调用
	TopicSubDeviceServiceInvokeReply = TopicPrefix + "+/+/thing/sub/service/invoke_reply" //设备->平台 设备服务调用响应
	TopicSubDeviceServiceInvoke      = TopicPrefix + "%s/%s/thing/sub/service/invoke"     //平台->设备 设备服务调用

)
