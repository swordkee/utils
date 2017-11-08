package utils

import (
	"errors"
	"math"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"bytes"
)

func IsEmpty(val interface{}) bool {
	v := reflect.ValueOf(val)
	valType := v.Kind()
	switch valType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Interface, reflect.Slice, reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func:
		// Check for empty slices and props
		if v.IsNil() {
			return true
		} else if valType == reflect.Slice || valType == reflect.Map {
			return v.Len() == 0
		}
	case reflect.Struct:
		fieldCount := v.NumField()
		for i := 0; i < fieldCount; i++ {
			field := v.Field(i)
			if field.IsValid() && !IsEmpty(field) {
				return false
			}
		}
		return true
	default:
		return false
	}
	return false
}

func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(s, "\\")
	return string(s[0: i+1])
}
func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

func CallFunc(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
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

func StringBuilder(str ...string) string {
	if len(str) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, v := range str {
		buffer.WriteString(v)
	}
	return buffer.String()
}
