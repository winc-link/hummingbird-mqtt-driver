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

type Sys struct {
	Ack bool `json:"ack"`
}

// PropertySet 属性设置
type PropertySet struct {
	MsgId   string                 `json:"msgId"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

func (p *PropertySet) Marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}

// PropertyQuery 设备属性查询
type PropertyQuery struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Params  []string `json:"params"`
}

func (p *PropertyQuery) Marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}

// PropertySetReply 设置设备属性响应
type PropertySetReply struct {
	Id   string                        `json:"id"`
	Data model.PropertySetResponseData `json:"data"`
}

func (p *PropertySetReply) Marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}

// PropertyQueryReply 设备属性实时查询响应
type PropertyQueryReply struct {
	Id     string                          `json:"id"`
	Params []model.PropertyGetResponseData `json:"params"`
}
