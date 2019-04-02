package utils

import (
	"testing"
)

func TestString(t *testing.T) {
	var i = "\n \t33 4 6999   66666    .677777\n"
	t.Log(i)
	t.Log(Bool(i))
	t.Log(BoolMust(i, true))
	t.Log(Bytes(i))
	t.Log(BytesMust(i))
	t.Log(Float32(i))
	t.Log(Float32Must(i, 6.666))
	t.Log(Float64(i))
	t.Log(Float64Must(i, 6.666))
	t.Log(Int(i))
	t.Log(IntMust(i, 6666))
	t.Log(Int8(i))
	t.Log(Int8Must(i, -5))
	t.Log(Int16(i))
	t.Log(Int16Must(i, -9))
	t.Log(Int32(i))
	t.Log(Int32Must(i, 6666))
	t.Log(Int64(i))
	t.Log(Int64Must(i, 6666))
	t.Log(Uint(i))
	t.Log(UintMust(i, 6666))
	t.Log(Uint8(i))
	t.Log(Uint8Must(i, 6))
	t.Log(Uint16(i))
	t.Log(Uint16Must(i, 6666))
	t.Log(Uint32(i))
	t.Log(Uint32Must(i, 6666))
	t.Log(Uint64(i))
	t.Log(Uint64Must(i, 6666))
	t.Log(String(i))
	t.Log(StringMust(i, "HHHH"))
}
