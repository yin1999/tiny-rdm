package convutil

import (
	"encoding/json"
	"reflect"
	"strconv"

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
		v := reflect.ValueOf(obj)
		switch v.Type().Kind() {
		case reflect.Map, reflect.Slice, reflect.Array:
			if b, err := json.Marshal(obj); err == nil {
				return string(b), true
			}
		case reflect.String:
			return obj.(string), true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(v.Int(), 10), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(v.Uint(), 10), true
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(v.Float(), 'f', -1, 64), true
		case reflect.Bool:
			return strconv.FormatBool(v.Bool()), true
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
