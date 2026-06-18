package jsonutil

import (
	json "github.com/bytedance/sonic"
)

func Decode[T any](data []byte, target T) *T {
	err := json.Unmarshal(data, &target)
	if err != nil {
		return nil
	}
	return &target
}

func RemoveNullValues(inputJSON string) string {
	// 解析 JSON 数据到 map[string]interface{}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(inputJSON), &data)
	if err != nil {
		return ""
	}

	// 遍历 map，删除值为 null 的键值对
	for key, value := range data {
		if value == nil {
			delete(data, key)
		}
	}

	// 将修改后的 map 转换回 JSON 字符串
	resultJSON, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(resultJSON)
}
