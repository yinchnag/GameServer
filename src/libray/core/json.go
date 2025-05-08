package core

import (
	jsoniter "github.com/json-iterator/go"
)

// 解编Json数据
func Unmarshal(data []byte, v interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}

// 编码Json数据
func Marshal(v interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(v)
}
