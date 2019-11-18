package utils

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}
func TestVariable(t *testing.T) {
	equal(t, true, IsEmpty(""))
	equal(t, true, IsEmpty(0))
	equal(t, true, IsEmpty(0.0))
	equal(t, true, IsEmpty(false))
	equal(t, false, IsEmpty([1]string{}))
	equal(t, true, IsEmpty([]int{}))

	tIsNumeric := IsNumeric("-0xaF")
	equal(t, true, tIsNumeric)
}
func TestArray(t *testing.T) {
	var s1 = make([]interface{}, 3)
	s1[0] = "a"
	s1[1] = "b"
	s1[2] = "c"
	tArrayChunk := ArrayChunk(s1, 2)
	equal(t, false, tArrayChunk)

	tInArray := InArray("a", s1)
	equal(t, true, tInArray)

	var m1 = make(map[interface{}]interface{}, 3)
	m1[1] = "a"
	m1["a"] = "b"
	m1[2.5] = 1
	tArrayKeyExists := ArrayKeyExists("a", m1)
	equal(t, true, tArrayKeyExists)

	tArrayUnshift := ArrayUnshift(&s1, "x", "y")
	equal(t, 5, tArrayUnshift)
	equal(t, []interface{}{"x", "y", "a", "b", "c"}, s1)

	tArrayPush := ArrayPush(&s1, "u", "v")
	equal(t, 7, tArrayPush)
	equal(t, []interface{}{"x", "y", "a", "b", "c", "u", "v"}, s1)

	tArrayPop := ArrayPop(&s1)
	equal(t, "v", tArrayPop)
	equal(t, []interface{}{"x", "y", "a", "b", "c", "u"}, s1)

	tArrayShift := ArrayShift(&s1)
	equal(t, "x", tArrayShift)
	equal(t, []interface{}{"y", "a", "b", "c", "u"}, s1)

	tarrayfill := ArrayFill(-3, 6, "aaa")
	equal(t, map[int]interface{}{-1: "aaa", 0: "aaa", 1: "aaa", 2: "aaa", -3: "aaa", -2: "aaa"}, tarrayfill)

	tarrayrand := ArrayRand([]interface{}{"a", "b", "c"})
	equal(t, 3, len(tarrayrand))

	var s2 = make([]interface{}, 3)
	s2[0] = "a"
	s2[1] = "b"
	s2[2] = "c"
	tarraypad := ArrayPad(s2, -5, "d")
	equal(t, []interface{}{"d", "d", "a", "b", "c"}, tarraypad)

	var s3 = make([]interface{}, 3, 3)
	s3[0] = "x"
	s3[1] = "y"
	s3[2] = "z"
	tarraycombine := ArrayCombine(s2, s3)
	equal(t, map[interface{}]interface{}{"a": "x", "b": "y", "c": "z"}, tarraycombine)
}
