package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// String : Conver "val" to a String
func String(val interface{}) (string, error) {
	switch ret := val.(type) {
	case string:
		return ret, nil
	case []byte:
		return string(ret), nil
	default:
		str := fmt.Sprintf("%+v", val)
		if val == nil || len(str) == 0 {
			return "", fmt.Errorf("conver.String(), the %+v is empty", val)
		}
		return str, nil
	}
}

// StringMust : Must Conver "val" to a String
func StringMust(val interface{}, def ...string) string {
	ret, err := String(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return strings.TrimSpace(ret)
}

// Bool : Conver "val" to a Bool
func Bool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}
	switch ret := val.(type) {
	case bool:
		return ret, nil
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return ret != 0, nil
	case []byte:
		return stringToBool(string(ret))
	case string:
		return stringToBool(ret)
	default:
		return false, converError(val, "bool")
	}
}

// BoolMust : Must Conver "val" to a Bool
func BoolMust(val interface{}, def ...bool) bool {
	ret, err := Bool(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Bytes : Conver "val" to []byte
func Bytes(val interface{}) ([]byte, error) {
	switch ret := val.(type) {
	case []byte:
		return ret, nil
	default:
		str, err := String(val)
		return []byte(str), err
	}
}

// BytesMust : Must Conver "val" to []byte
func BytesMust(val interface{}, def ...[]byte) []byte {
	ret, err := Bytes(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float32 : Conver "val" to a Float32
func Float32(val interface{}) (float32, error) {
	switch ret := val.(type) {
	case float32:
		return ret, nil
	case int:
		return float32(ret), nil
	case int8:
		return float32(ret), nil
	case int16:
		return float32(ret), nil
	case int32:
		return float32(ret), nil
	case int64:
		return float32(ret), nil
	case uint:
		return float32(ret), nil
	case uint8:
		return float32(ret), nil
	case uint16:
		return float32(ret), nil
	case uint32:
		return float32(ret), nil
	case uint64:
		return float32(ret), nil
	case float64:
		return float32(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		f, err := strconv.ParseFloat(str, 32)
		return float32(f), err
	}
}

// Float32Must : Must Conver "val" to Float32
func Float32Must(val interface{}, def ...float32) float32 {
	ret, err := Float32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float64 : Conver "val" to a Float64
func Float64(val interface{}) (float64, error) {
	switch ret := val.(type) {
	case float64:
		return ret, nil
	case int:
		return float64(ret), nil
	case int8:
		return float64(ret), nil
	case int16:
		return float64(ret), nil
	case int32:
		return float64(ret), nil
	case int64:
		return float64(ret), nil
	case uint:
		return float64(ret), nil
	case uint8:
		return float64(ret), nil
	case uint16:
		return float64(ret), nil
	case uint32:
		return float64(ret), nil
	case uint64:
		return float64(ret), nil
	case float32:
		return float64(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		return strconv.ParseFloat(str, 64)
	}
}

// Float64Must : Must Conver "val" to Float64
func Float64Must(val interface{}, def ...float64) float64 {
	ret, err := Float64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Uint : Conver "val" to a rounded Uint
func Uint(val interface{}) (uint, error) {
	i, err := Uint64(val)
	if err != nil {
		return 0, err
	}
	return uint(i), err
}

// UintMust : Must Conver "val" to a rounded Uint
func UintMust(val interface{}, def ...uint) uint {
	ret, err := Uint(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Uint8 : Conver "val" to a rounded uint8
func Uint8(val interface{}) (uint8, error) {
	i, err := Uint64(val)
	if err != nil {
		return 0, err
	}
	return uint8(i), err
}

// Uint8Must : Must Conver "val" to a rounded uint8
func Uint8Must(val interface{}, def ...uint8) uint8 {
	ret, err := Uint8(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Uint16 : Conver "val" to a rounded uint16
func Uint16(val interface{}) (uint16, error) {
	i, err := Uint64(val)
	if err != nil {
		return 0, err
	}
	return uint16(i), err
}

// Uint16Must : Must Conver "val" to a rounded uint16
func Uint16Must(val interface{}, def ...uint16) uint16 {
	ret, err := Uint16(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Uint32 : Conver "val" to a rounded uint32
func Uint32(val interface{}) (uint32, error) {
	i, err := Uint64(val)
	if err != nil {
		return 0, err
	}
	return uint32(i), err
}

// Uint32Must : Must Conver "val" to a rounded uint32
func Uint32Must(val interface{}, def ...uint32) uint32 {
	ret, err := Uint32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Uint64 : Conver "val" to a rounded Uint64
func Uint64(val interface{}) (uint64, error) {
	str, err := String(val)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(str, 10, 64)
}

// Uint64Must : Must Conver "val" to a rounded uint64
func Uint64Must(val interface{}, def ...uint64) uint64 {
	ret, err := Uint64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int : Conver "val" to a rounded Int
func Int(val interface{}) (int, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}
	return int(i), err
}

// IntMust : Must Conver "val" to a rounded Int
func IntMust(val interface{}, def ...int) int {
	ret, err := Int(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int8 : Conver "val" to a rounded Int8
func Int8(val interface{}) (int8, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}
	return int8(i), err
}

// Int8Must : Must Conver "val" to a rounded Int8
func Int8Must(val interface{}, def ...int8) int8 {
	ret, err := Int8(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int16 : Conver "val" to a rounded Int16
func Int16(val interface{}) (int16, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}
	return int16(i), err
}

// Int16Must : Must Conver "val" to a rounded Int16
func Int16Must(val interface{}, def ...int16) int16 {
	ret, err := Int16(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int32 : Conver "val" to a rounded Int32
func Int32(val interface{}) (int32, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}
	return int32(i), err
}

// Int32Must : Must Conver "val" to a rounded Int32
func Int32Must(val interface{}, def ...int32) int32 {
	ret, err := Int32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int64 : Conver "val" to a rounded Int64
func Int64(val interface{}) (int64, error) {
	str, err := String(val)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 10, 64)
}

// Int64Must : Must Conver "val" to a rounded Int64
func Int64Must(val interface{}, def ...int64) int64 {
	ret, err := Int64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}
func converError(val interface{}, t string) error {
	return fmt.Errorf("conver error, the %T{%v} can not conver to a %v", val, val, t)
}

func stringToBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "ok", "OK", "yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "":
		return false, nil
	}
	return false, converError(val, "bool")
}
