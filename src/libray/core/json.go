package core

import (
	json "github.com/bytedance/sonic"
)

// 解编Json数据
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// 编码Json数据
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
