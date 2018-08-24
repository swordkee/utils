package utils

import (
	"gonum.org/v1/gonum/floats"
	"reflect"
)

//func init() {
//	gob.Register([]M{})
//	gob.Register(M{})
//}
//
//type M map[string]interface{}

func InArray(val interface{}, array interface{}) (exists bool) {
	if InArrayIdx(val, array) != -1 {
		return true
	}
	return false
}
func InArrayIdx(val interface{}, array interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}

	return
}
func InArrayBetween(val float64, array []float64) (exists bool) {
	if val == 0 {
		return false
	}
	min := floats.Min(array)
	max := floats.Max(array)
	if min > 0 && max > 0 && min <= val && val <= max ||
		InArray(val, array) {
		return true
	}
	return false
}
func ArrayColumn(d map[int]map[string]string, column_key string, index_key string) map[int]map[string]string {
	nd := make(map[int]map[string]string)
	for k, v := range d {
		for e, q := range v {
			if e == index_key {
				nd[k][index_key] = q
			}
		}
	}
	return nd
}
func ArrayKeys(s string, d map[string]interface{}) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}
func ArrayValues(d map[string]interface{}) []interface{} {
	nd := make([]interface{}, 0)
	for _, v := range d {
		if v != "" {
			nd = append(nd, v)
		}
	}
	return nd
}
func ArrayFlip(d map[string]string) map[string]string {
	nd := make(map[string]string, 0)
	for k, v := range d {
		if v != "" {
			nd[v] = k
		}
	}
	return nd
}
func ArrayMerge(s ...[]interface{}) (slice []interface{}) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := ArrayMerge(s[1:]...) //...将数组元素打散
		slice = make([]interface{}, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}

	return
}

// convert map to struct
/*
func (m M) MapToStruct(s interface{}) {
	v := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		key := f.Name
		scnKey := Strings(key).SnakeCasedName()
		tag := f.Tag
		fieldName := tag.Get("field")
		vf := v.Field(i)
		doRes := false
		if fieldName != "" {
			if val, ok := m[fieldName]; ok {
				vv := reflect.ValueOf(val)
				if vf.Type().Kind().String() != vv.Type().Kind().String() {
					if vf.Type().Kind().String() == "bool" {
						if vv.Type().Kind().String() == "int64" && vv.Int() > 0 {
							vf.SetBool(true)
						}

						if vv.Type().Kind().String() == "string" && vv.String() != "" {
							ii, _ := strconv.ParseInt(vv.String(), 10, 64)
							if ii > 0 {
								vf.SetBool(true)
							}
						}
					}
				} else {
					vf.Set(vv)
				}

				doRes = true
			}
		}

		if !doRes {
			if _, ok := m[key]; ok {
				vf.Set(reflect.ValueOf(m[key]))
			}
		}

		if !doRes {
			if _, ok := m[scnKey]; ok {
				vf.Set(reflect.ValueOf(m[scnKey]))
			}
		}
	}
}*/
