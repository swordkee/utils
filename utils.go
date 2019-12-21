package utils

import (
	"encoding/binary"
	"errors"
	"github.com/json-iterator/go"
	"math"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
)

func IsEmpty(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// is_numeric()
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.TrimSpace(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

func CallFunc(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

//CallFunc(funcs, "foo")
//CallFunc(funcs, "bar", 1, 2, 3)

//四舍五入 取精度
func ToFixed(f float64, places int) float64 {
	if f == 0 {
		return 0
	}
	shift := math.Pow(10, float64(places))
	fv := 0.0000000001 + f //对浮点数产生.xxx999999999 计算不准进行处理
	return math.Floor(fv*shift+.5) / shift
}

//笛卡尔乘积
func Iter(params ...[]interface{}) chan []interface{} {
	// create channel
	c := make(chan []interface{})
	// create waitgroup
	var wg sync.WaitGroup
	// call iterator
	wg.Add(1)
	iterate(&wg, c, []interface{}{}, params...)
	// call channel-closing go-func
	go func() { wg.Wait(); close(c) }()
	// return channel
	return c
}

// private, recursive Iteration-Function
func iterate(wg *sync.WaitGroup, channel chan []interface{}, result []interface{}, params ...[]interface{}) {
	// dec WaitGroup when finished
	defer wg.Done()
	// no more params left?
	if len(params) == 0 {
		// send result to channel
		channel <- result
		return
	}
	// shift first param
	p, params := params[0], params[1:]
	// iterate over it
	for i := 0; i < len(p); i++ {
		// inc WaitGroup
		wg.Add(1)
		// call self with remaining params
		go iterate(wg, channel, append(result, p[i]), params...)
	}
}

//三元
// On(a ? b : c)
func On(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
func MapToStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
}
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

func JsonEncode(val interface{}) string {
	var jso = jsoniter.ConfigCompatibleWithStandardLibrary
	ret, err := jso.Marshal(val)
	if err != nil {
		return ""
	}
	return string(ret)
}

func JsonDecode(data string, val interface{}) error {
	var jso = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := jso.Unmarshal([]byte(data), &val); err != nil {
		return err
	}
	return nil
}

// Exit exit()
func Exit(status int) {
	os.Exit(status)
}

// Die die()
func Die(status int) {
	os.Exit(status)
}

// IP2long ip2long()
// IPv4
func IP2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// Long2ip long2ip()
// IPv4
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}
