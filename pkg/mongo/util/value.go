package util

import "reflect"

// IsZeroValue checks if the given value is the zero value for its type
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return val == 0
	case uint, uint8, uint16, uint32, uint64:
		return val == 0
	case float32, float64:
		return val == 0
	case bool:
		return val
	default:
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Func, reflect.Chan:
			return rv.IsNil()
		case reflect.Slice:
			return rv.IsNil() || rv.Len() == 0
		case reflect.Struct:
			return reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
		default:
			return reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
		}
	}
}
