package deviceonline

import (
	"github.com/winc-link/hummingbird-mqtt-driver/config"
	constants "github.com/winc-link/hummingbird-mqtt-driver/constant"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"sync"
	"time"
)

// DeviceStatusTimeController实现设备在线状态控制
var sd *service.DriverService

func InitDeviceStatusTimeController(s *service.DriverService) {
	sd = s
}

var devicesMap sync.Map

type deviceTimeController struct {
	deviceId       string
	lastReportTime int64
	onLineType     string
	interval       int
	timer          *time.Timer
}

func newDeviceTimeController(deviceId string, interval int, o string) *deviceTimeController {
	var dt deviceTimeController
	dt.deviceId = deviceId
	dt.lastReportTime = 0
	dt.interval = interval
	dt.onLineType = o
	dt.timer = time.NewTimer(time.Duration(interval) * time.Minute)
	return &dt
}

// AddDeviceStatusTimeController
// 在DeviceNotify回调方法中如果DeviceNotifyType是add需要调用此方法
// 驱动启动成功需要调用次方法
// 通过SDK创建设备时需要手动调用次方法
func AddDeviceStatusTimeController(device model.Device) {
	interval := config.GetConfig().Interval
	if interval >= 1 {
		//判断devicesMap是否存在这个设备
		_, ok := devicesMap.Load(device.Id)
		if !ok {
			dt := newDeviceTimeController(device.Id, interval, constants.CustomOnline)
			devicesMap.Store(dt.deviceId, dt)
			go runLoopCheck(dt)
		}
	}
}

// DelDeviceStatusTimeController
// 在DeviceNotify回调方法中如果DeviceNotifyType是delete需要调用此方法
func DelDeviceStatusTimeController(deviceId string) {
	dt, ok := devicesMap.Load(deviceId)
	if ok {
		if dt.(*deviceTimeController).timer != nil {
			dt.(*deviceTimeController).timer.Stop()
			devicesMap.Delete(deviceId)
		}
	}
}

// UpdateDeviceLastReportTime
// 在调用SDK中PropertyReport方法之前需要调用此方法记录设备最后一次上报时间
func UpdateDeviceLastReportTime(deviceId string) {
	dt, ok := devicesMap.Load(deviceId)
	if ok {
		dt.(*deviceTimeController).lastReportTime = time.Now().Unix()
		devicesMap.Store(deviceId, dt)
	}
}

// runLoopCheck 设备轮训定时器
func runLoopCheck(dt *deviceTimeController) {
	go func(t *time.Timer) {
		for {
			<-t.C
			deviceC, ok := devicesMap.Load(dt.deviceId)
			if ok {
				//在线离线根据上报时间做判断
				deviceTC, ok := deviceC.(*deviceTimeController)
				if ok {
					if deviceTC.lastReportTime+int64(deviceTC.interval)*60 < time.Now().Unix() {
						device, ok := sd.GetDeviceById(deviceTC.deviceId)
						if ok {
							if device.Status == commons.DeviceOnline {
								sd.GetLogger().Info("更新设备状态为离线：deviceId:", deviceTC.deviceId)
								err := sd.Offline(deviceTC.deviceId)
								if err != nil {
									sd.GetLogger().Error(err)
								}
							}
						}
					} else {
						device, ok := sd.GetDeviceById(deviceTC.deviceId)
						if ok {
							if device.Status == commons.DeviceOffline {
								sd.GetLogger().Info("更新设备状态为在线：deviceId:", deviceTC.deviceId)
								err := sd.Online(deviceTC.deviceId)
								if err != nil {
									sd.GetLogger().Error(err)
								}
							}
						}
					}
					//sd.GetLogger().Info(time.Now().Format("2006-01-02 15:04:05"), "开始重置时间，时间间隔为：", deviceTC.interval)
					t.Reset(time.Duration(deviceTC.interval) * time.Minute)
				}
			} else {
				t.Stop()
				return
			}

		}
	}(dt.timer)
}
