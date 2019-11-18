package utils

import (
	"math"
	"math/rand"
	"reflect"
	"time"
)

//in_array()
func InArray(val interface{}, array interface{}) (exists bool) {
	if InArrayIdx(val, array) != -1 {
		return true
	}
	return false
}

func InArrayIdx(val interface{}, array interface{}) (index int) {
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
func InArrayBetween(val float64, array []float64) (exists bool) {
	if val == 0 {
		return false
	}
	min := Min(array)
	max := Max(array)
	if min > 0 && max > 0 && min <= val && val <= max ||
		InArray(val, array) {
		return true
	}
	return false
}

func ArrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

//array_keys
func ArrayKeys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key, _ := range elements {
		keys[i] = key
		i++
	}
	return keys
}

//array_values
func ArrayValues(elements map[interface{}]interface{}) []interface{} {
	nd := make([]interface{}, len(elements))
	for _, v := range elements {
		if v != "" {
			nd = append(nd, v)
		}
	}
	return nd
}

//array_flip
func ArrayFlip(m map[interface{}]interface{}) map[interface{}]interface{} {
	nd := make(map[interface{}]interface{})
	for k, v := range m {
		if v != "" {
			nd[v] = k
		}
	}
	return nd
}

//array_merge()
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

// array_fill()
func ArrayFill(startIndex int, num uint, value interface{}) map[int]interface{} {
	m := make(map[int]interface{})
	var i uint
	for i = 0; i < num; i++ {
		m[startIndex] = value
		startIndex++
	}
	return m
}

// array_chunk()
func ArrayChunk(s []interface{}, size int) [][]interface{} {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

// array_pad()
func ArrayPad(s []interface{}, size int, val interface{}) []interface{} {
	if size == 0 || (size > 0 && size < len(s)) || (size < 0 && size > -len(s)) {
		return s
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(s)
	tmp := make([]interface{}, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		return append(s, tmp...)
	} else {
		return append(tmp, s...)
	}
}

// array_slice()
func ArraySlice(s []interface{}, offset, length uint) []interface{} {
	if offset > uint(len(s)) {
		panic("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(s)) {
		return s[offset:end]
	}
	return s[offset:]
}

// array_rand()
func ArrayRand(elements []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]interface{}, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// array_push()
// Push one or more elements onto the end of slice
func ArrayPush(s *[]interface{}, elements ...interface{}) int {
	*s = append(*s, elements...)
	return len(*s)
}

// array_pop()
// Pop the element off the end of slice
func ArrayPop(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// array_unshift()
// Prepend one or more elements to the beginning of a slice
func ArrayUnshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// array_shift()
// Shift an element off the beginning of slice
func ArrayShift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// array_key_exists()
func ArrayKeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}

// array_combine()
func ArrayCombine(s1, s2 []interface{}) map[interface{}]interface{} {
	if len(s1) != len(s2) {
		panic("the number of elements for each slice isn't equal")
	}
	m := make(map[interface{}]interface{}, len(s1))
	for i, v := range s1 {
		m[v] = s2[i]
	}
	return m
}

// array_reverse()
func ArrayReverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
