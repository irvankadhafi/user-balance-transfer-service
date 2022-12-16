package utils

import "encoding/json"

// InterfaceBytesToType will transform cached value that get from the redis to any types
func InterfaceBytesToType[T any](i interface{}) (out T) {
	if i == nil {
		return
	}
	bt := i.([]byte)

	_ = json.Unmarshal(bt, &out)
	return
}
