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
	"crypto/tls"
	"net"
	"net/http"
	"os"

	"github.com/winc-link/hummingbird-sdk-go/service"

	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/pkg/pidfile"
	"github.com/DrmagicE/gmqtt/server"

	_ "github.com/DrmagicE/gmqtt/persistence"
	_ "github.com/DrmagicE/gmqtt/plugin/prometheus"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
)

var GlobalDriverService *service.DriverService

type MQTTServer struct {
}

func NewMQTTService(sd *service.DriverService) *MQTTServer {
	GlobalDriverService = sd
	return &MQTTServer{}
}

func (m *MQTTServer) Start() {
	c := config.DefaultConfig()
	c.Listeners = DefaultListeners
	if c.PidFile != "" {
		pid, err := pidfile.New(c.PidFile)
		if err != nil {
			panic(err)
		}
		defer pid.Remove()
	}

	tcpListeners, websockets, err := GetListeners(c)
	if err != nil {
		GlobalDriverService.GetLogger().Error(err.Error())
		os.Exit(1)
	}

	//l, err := c.GetLogger(c.Log)
	//if err != nil {
	//	GlobalDriverService.GetLogger().Error(err.Error())
	//	os.Exit(1)
	//}

	//设置Hooks
	hooks := server.Hooks{
		OnAccept:            OnAccept,
		OnStop:              OnStop,
		OnSubscribe:         OnSubscribe,
		OnSubscribed:        OnSubscribed,
		OnUnsubscribe:       OnUnsubscribe,
		OnUnsubscribed:      OnUnsubscribed,
		OnMsgArrived:        OnMsgArrived,
		OnBasicAuth:         OnBasicAuth,
		OnEnhancedAuth:      OnEnhancedAuth,
		OnReAuth:            OnReAuth,
		OnConnected:         OnConnected,
		OnSessionCreated:    OnSessionCreated,
		OnSessionResumed:    OnSessionResumed,
		OnSessionTerminated: OnSessionTerminated,
		OnDelivered:         OnDelivered,
		OnClosed:            OnClosed,
		OnMsgDropped:        OnMsgDropped,
		OnWillPublish:       OnWillPublish,
		OnWillPublished:     OnWillPublished,
	}
	s := server.New(
		server.WithConfig(c),
		server.WithTCPListener(tcpListeners...),
		server.WithWebsocketServer(websockets...),
		server.WithHook(hooks),
		//server.WithLogger(l),
	)

	err = s.Init()
	if err != nil {
		GlobalDriverService.GetLogger().Error(err.Error())
		os.Exit(1)
	}

	//启动server
	err = s.Run()
	if err != nil {
		GlobalDriverService.GetLogger().Error(err.Error())
		os.Exit(1)
	}
}

func GetListeners(c config.Config) (tcpListeners []net.Listener, websockets []*server.WsServer, err error) {
	for _, v := range c.Listeners {
		var ln net.Listener
		if v.Websocket != nil {
			ws := &server.WsServer{
				Server: &http.Server{Addr: v.Address},
				Path:   v.Websocket.Path,
			}
			if v.TLSOptions != nil {
				ws.KeyFile = v.Key
				ws.CertFile = v.Cert
			}
			websockets = append(websockets, ws)
			continue
		}
		if v.TLSOptions != nil {
			var cert tls.Certificate
			cert, err = tls.LoadX509KeyPair(v.Cert, v.Key)
			if err != nil {
				return
			}
			ln, err = tls.Listen("tcp", v.Address, &tls.Config{
				Certificates: []tls.Certificate{cert},
			})
		} else {
			ln, err = net.Listen("tcp", v.Address)
		}
		tcpListeners = append(tcpListeners, ln)
	}
	return
}
