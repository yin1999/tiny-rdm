package convutil

import (
	"encoding/json"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackConvert struct{}

func (MsgpackConvert) Enable() bool {
	return true
}

func (c MsgpackConvert) Encode(str string) (string, bool) {
	var obj any
	if err := json.Unmarshal([]byte(str), &obj); err == nil {
		obj = c.TryFloatToInt(obj)
		if b, err := msgpack.Marshal(obj); err == nil {
			return string(b), true
		}
	}

	if b, err := msgpack.Marshal(str); err == nil {
		return string(b), true
	}

	return str, false
}

func (MsgpackConvert) Decode(str string) (string, bool) {
	var obj any
	if err := msgpack.Unmarshal([]byte(str), &obj); err == nil {
		t := reflect.TypeOf(obj)
		switch t.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			if b, err := json.Marshal(obj); err == nil {
				return string(b), true
			}
		case reflect.String:
			return obj.(string), true
		}
	}

	return str, false
}

func (c MsgpackConvert) TryFloatToInt(input any) any {
	switch val := input.(type) {
	case map[string]any:
		for k, v := range val {
			val[k] = c.TryFloatToInt(v)
		}
		return val
	case []any:
		for i, v := range val {
			val[i] = c.TryFloatToInt(v)
		}
		return val
	case float64:
		if val == float64(int(val)) {
			return int(val)
		}
		return val
	default:
		return val
	}
}
