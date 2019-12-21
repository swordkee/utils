package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// callFn if args[0] == func, run it
func callFn(f interface{}) interface{} {
	if v, ok := f.(func() interface{}); ok {
		return v()
	}
	if v, ok := f.(func()); ok {
		v()
		return nil
	}
	return f
}

// TestFnTime run func use time
func TestFnTime(f interface{}) string {
	start := time.Now()
	callFn(f)
	end := time.Now()
	vf := reflect.ValueOf(f)
	str := fmt.Sprintf("[%s] runtime: %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
	fmt.Println(str)
	return str
}

// If - (a ? b : c) Or (a && b)
func If(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	var trueVal = args[1]
	var falseVal interface{}
	if len(args) > 2 {
		falseVal = args[2]
	} else {
		falseVal = nil
	}
	if condition == nil {
		return callFn(falseVal)
	} else if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(falseVal)
		}
	} else if v, ok := condition.(int); ok {
		if v == 0 {
			return callFn(falseVal)
		}
	} else if v, ok := condition.(string); ok {
		if v != "" && v != "0" && v != "false" {
			return callFn(trueVal)
		}
		return callFn(falseVal)
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return callFn(trueVal)
}

// Or - (a || b)
func Or(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	if condition == nil {
		return callFn(args[1])
	}
	if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(args[1])
		}
	} else if v, ok := condition.(int); ok {
		if v == 0 {
			return callFn(args[1])
		}
	} else if v, ok := condition.(string); ok {
		if v != "" && v != "0" && v != "false" {
			return condition
		}
		return callFn(args[1])
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return condition
}

func In(val interface{}, array interface{}) (exists bool) {
	if InIdx(val, array) != -1 {
		return true
	}
	return false
}

func InIdx(val interface{}, array interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	case reflect.Map:
		s := reflect.ValueOf(array)
		if s.MapIndex(reflect.ValueOf(val)).IsValid() {
			index = 0
			return
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}
	return
}
