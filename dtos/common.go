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

import "encoding/json"

type CommonResponse struct {
	Id           string `json:"id"`
	Code         int    `json:"code"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Version      string `json:"version"`
}

func (t *CommonResponse) Marshal() []byte {
	b, _ := json.Marshal(t)
	return b
}

type IntOrFloatSpecs struct {
	Min      string `json:"min"`
	Max      string `json:"max"`
	Step     string `json:"step"`
	Unit     string `json:"unit"`
	UnitName string `json:"unitName"`
}

type TextSpecs struct {
	Length string `json:"length"`
}
