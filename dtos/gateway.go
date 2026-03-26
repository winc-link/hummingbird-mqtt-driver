package dtos

import (
	"encoding/json"

	"github.com/spf13/cast"
)

type WifiInfo struct {
	IsSupport  bool   `json:"is_support"`
	WifiName   string `json:"wifi_name"`
	WifiSignal string `json:"wifi_signal"`
}
type OperatorInfo struct {
	IsSupport bool   `json:"is_support"`
	Operator  string `json:"operator"`
	Iccid     string `json:"iccid"`
	Imei      string `json:"imei"`
	Signal    string `json:"signal"`
}
type Monitor struct {
	MemoryTotal string  `json:"memory_total"`
	MemoryUsed  string  `json:"memory_used"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskTotal   string  `json:"disk_total"`
	DiskUsed    string  `json:"disk_used"`
	DiskUsage   float64 `json:"disk_usage"`
	CpuTotal    string  `json:"cpu_total"`
	CpuUsage    float64 `json:"cpu_usage"`
}

type GatewayInfo struct {
	GatewaySn          string              `json:"gateway_sn"`
	Name               string              `json:"gateway_name"`
	Secert             string              `json:"gateway_secert"` //网关密钥
	Model              string              `json:"gateway_model"`  //网关型号
	TimeZone           string              `json:"gateway_timezone"`
	RunTime            string              `json:"gateway_runtime"`
	NumDevice          int                 `json:"num_device"` //子设备数量
	NumPoint           int                 `json:"num_point"`  //点位数
	NumPointStatus     []map[string]any    `json:"num_point_status"`
	NumDeviceProtocols []map[string]any    `json:"num_device_protocols"`
	IpAddress          []map[string]string `json:"ip_address"` //设备的ip
	WifiInfo           WifiInfo            `json:"wifi_info"`
	OperatorInfo       OperatorInfo        `json:"operator_info"`
	//Lat              string                 `json:"lat"`            //设备位置
	//Lon              string                 `json:"lon"`            //设备位置
	//Location         string                 `json:"location"`       //设备位置
	// FrpIp        string   `json:"frp_ip"`         //IP地址
	// FrpPort      string   `json:"frp_port"`       //端口
	Monitor Monitor `json:"monitor"`
}

type GatewayReportData struct {
	Params    GatewayInfo `json:"params"`
	Timestamp int64       `json:"Timestamp"`
	MsgId     string      `json:"msg_id"`
}

func (g GatewayInfo) TransFormRedis() map[string]string {
	res := make(map[string]string)
	res["gateway_sn"] = g.GatewaySn
	res["gateway_name"] = g.Name
	res["gateway_secert"] = g.Secert
	res["gateway_model"] = g.Model
	res["gateway_timezone"] = g.TimeZone
	res["gateway_runtime"] = g.RunTime

	ip_address, _ := json.Marshal(g.IpAddress)
	res["ip_address"] = (string)(ip_address)
	num_device_protocols, _ := json.Marshal(g.NumDeviceProtocols)
	res["num_device_protocols"] = (string)(num_device_protocols)
	res["num_device"] = cast.ToString(g.NumDevice)
	operator_info, _ := json.Marshal(g.OperatorInfo)
	res["operator_info"] = string(operator_info)
	wifi_info, _ := json.Marshal(g.WifiInfo)
	res["wifi_info"] = string(wifi_info)
	res["num_point"] = cast.ToString(g.NumPoint)
	num_point_status, _ := json.Marshal(g.NumPointStatus)
	res["num_point_status"] = string(num_point_status)
	monitor, _ := json.Marshal(g.Monitor)
	res["monitor"] = string(monitor)

	return res

}

func (g GatewayInfo) TransFormOrm() map[string]any {
	res := make(map[string]any)
	// res["gateway_sn"] = g.Id
	// res["gateway_name"] = g.Name
	res["gateway_secert"] = g.Secert
	res["gateway_model"] = g.Model
	res["gateway_timezone"] = g.TimeZone
	res["gateway_runtime"] = g.RunTime

	ip_address, _ := json.Marshal(g.IpAddress)
	res["ip_address"] = (string)(ip_address)
	num_device_protocols, _ := json.Marshal(g.NumDeviceProtocols)
	res["num_device_protocols"] = (string)(num_device_protocols)
	res["num_device"] = (g.NumDevice)
	operator_info, _ := json.Marshal(g.OperatorInfo)
	res["operator_info"] = string(operator_info)
	wifi_info, _ := json.Marshal(g.WifiInfo)
	res["wifi_info"] = string(wifi_info)
	res["num_point"] = (g.NumPoint)
	num_point_status, _ := json.Marshal(g.NumPointStatus)
	res["num_point_status"] = string(num_point_status)
	monitor, _ := json.Marshal(g.Monitor)
	res["monitor"] = string(monitor)

	return res
}
