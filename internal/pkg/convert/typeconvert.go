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

package convert

import (
	"fmt"
	"strconv"
)

func StringToInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	return i, err
}

func StringToFloat64(value string) (float64, error) {
	i, err := strconv.ParseFloat(value, 64)
	return i, err
}

func GetInterfaceToString(value interface{}) (string, error) {
	var t2 string
	switch value.(type) {
	case string:
		t2 = value.(string)
	default:
		return "", fmt.Errorf("convert error")
	}
	return t2, nil
}

func GetInterfaceToInt(t1 interface{}) (int, error) {
	var t2 int
	switch t1.(type) {
	case int:
		t2 = int(t1.(int))
		break
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	default:
		return 0, fmt.Errorf("convert error")
	}
	return t2, nil
}

func GetInterfaceToFloat64(t1 interface{}) (float64, error) {
	var t2 float64
	switch t1.(type) {
	case float64:
		t2 = t1.(float64)
		break
	case float32:
		t2 = float64(t1.(float32))
		break
	case int:
		t2 = float64(t1.(int))
		break
	case uint:
		t2 = float64(t1.(uint))
		break
	case int8:
		t2 = float64(t1.(int8))
		break
	case uint8:
		t2 = float64(t1.(uint8))
		break
	case int16:
		t2 = float64(t1.(int16))
		break
	case uint16:
		t2 = float64(t1.(uint16))
		break
	case int32:
		t2 = float64(t1.(int32))
		break
	case uint32:
		t2 = float64(t1.(uint32))
		break
	case int64:
		t2 = float64(t1.(int64))
		break
	case uint64:
		t2 = float64(t1.(uint64))
		break
	default:
		return 0, fmt.Errorf("convert error")
	}
	return t2, nil
}
