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

package dtos

import (
	"encoding/json"
	"github.com/winc-link/hummingbird-sdk-go/model"
)

// ServiceInvoke 服务调用
type ServiceInvoke struct {
	Id      string              `json:"id"`
	Version string              `json:"version"`
	Sys     Sys                 `json:"sys"`
	Params  model.ServiceDataIn `json:"params"`
}

func (s *ServiceInvoke) Marshal() []byte {
	b, _ := json.Marshal(s)
	return b
}

// ServiceInvokeReply 服务调用响应
type ServiceInvokeReply struct {
	Id     string               `json:"id"`
	Params model.ServiceDataOut `json:"params"`
}
