package utils

import (
	"reflect"
	"strconv"
	"testing"
)

// CheckInvariants checks that the invariants for s.data hold.
func CheckInvariants(t *testing.T, msg string, s *BitSet) {
	len := len(s.data)
	cap := cap(s.data)
	data := s.data[:cap]
	m := "Invariant for "
	if len > 0 && data[len-1] == 0 {
		t.Errorf("%s%s: data = %v, data[%d] = 0; want non-zero", m, msg, data, len-1)
	}
	for i := len; i < cap; i++ {
		if data[i] != 0 {
			t.Errorf("%s%s: data = %v, data[%d] = %#x; want 0", m, msg, data, i, data[i])
			break
		}
	}
}

// Panics tells if function f panics with parameters p.
func PanicsByBit(f interface{}, p ...interface{}) bool {
	fv := reflect.ValueOf(f)
	ft := reflect.TypeOf(f)
	if ft.NumIn() != len(p) {
		panic("wrong argument count")
	}
	pv := make([]reflect.Value, len(p))
	for i, v := range p {
		if reflect.TypeOf(v) != ft.In(i) {
			panic("wrong argument type")
		}
		pv[i] = reflect.ValueOf(v)
	}
	return call(fv, pv)
}

func call(fv reflect.Value, pv []reflect.Value) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	fv.Call(pv)
	return
}

func TestNewBitSet(t *testing.T) {
	for _, s := range []*BitSet{
		NewBitSet(),
		NewBitSet(-1),
		NewBitSet(1),
		NewBitSet(1, 1),
		NewBitSet(65),
		NewBitSet(1, 2, 3),
		NewBitSet(100, 200, 300),
	} {
		CheckInvariants(t, "NewBitSet", s)
	}
}

func TestContains(t *testing.T) {
	for _, x := range []struct {
		s        *BitSet
		n        int
		contains bool
	}{
		{NewBitSet(), -1, false},
		{NewBitSet(), 1, false},
		{NewBitSet(), 100, false},
		{NewBitSet(-1), 1, false},
		{NewBitSet(-1), -1, false},
		{NewBitSet(0), 0, true},
		{NewBitSet(1), 0, false},
		{NewBitSet(1), 1, true},
		{NewBitSet(1), 100, false},
		{NewBitSet(65), 0, false},
		{NewBitSet(65), 1, false},
		{NewBitSet(65), 65, true},
		{NewBitSet(65), 100, false},

		{NewBitSet(1, 2, 3), 0, false},
		{NewBitSet(1, 2, 3), 1, true},
		{NewBitSet(1, 2, 3), 2, true},
		{NewBitSet(1, 2, 3), 3, true},
		{NewBitSet(1, 2, 3), 4, false},

		{NewBitSet(100, 200, 300), 0, false},
		{NewBitSet(100, 200, 300), 100, true},
		{NewBitSet(100, 200, 300), 200, true},
		{NewBitSet(100, 200, 300), 300, true},
		{NewBitSet(100, 200, 300), 400, false},
	} {
		s, n := x.s, x.n
		contains := s.Contains(n)
		if contains != x.contains {
			t.Errorf("%v.Contains(%d) = %t; want %t", s, n, contains, x.contains)
		}
	}
}

func TestCmp(t *testing.T) {
	Zero, One := NewBitSet(), NewBitSet(1)
	for _, x := range []struct {
		s1, s2        *BitSet
		equal, subset bool
	}{
		{Zero, Zero, true, true},
		{One, One, true, true},
		{NewBitSet(), NewBitSet(), true, true},
		{NewBitSet(1), NewBitSet(1), true, true},
		{NewBitSet(64), NewBitSet(64), true, true},
		{NewBitSet(65), NewBitSet(65), true, true},
		{NewBitSet(1, 2, 3), NewBitSet(1, 2, 3), true, true},
		{NewBitSet(100, 200, 300), NewBitSet(100, 200, 300), true, true},

		{NewBitSet(), NewBitSet(1), false, true},
		{NewBitSet(1), NewBitSet(), false, false},
		{NewBitSet(1), NewBitSet(2), false, false},
		{NewBitSet(), NewBitSet(65), false, true},
		{NewBitSet(65), NewBitSet(), false, false},
		{NewBitSet(1), NewBitSet(65), false, false},
		{NewBitSet(1, 2, 3), NewBitSet(100, 200, 300), false, false},

		{NewBitSet(1), NewBitSet(1, 2, 3), false, true},
		{NewBitSet(1, 2, 3), NewBitSet(1), false, false},
		{NewBitSet(100), NewBitSet(100, 200, 300), false, true},
		{NewBitSet(100, 200, 300), NewBitSet(100), false, false},
	} {
		s1, s2 := x.s1, x.s2
		equal := s1.Equal(s2)
		if equal != x.equal {
			t.Errorf("%v.Equal(%v) = %t; want %t", s1, s2, equal, x.equal)
		}
		subset := s1.Subset(s2)
		if subset != x.subset {
			t.Errorf("%v.Subset(%v) = %t; want %t", s1, s2, subset, x.subset)
		}
	}
}

func TestMax(t *testing.T) {
	for _, x := range []struct {
		s   *BitSet
		max int
	}{
		{NewBitSet(0), 0},
		{NewBitSet(65), 65},
		{NewBitSet(1, 2, 3), 3},
		{NewBitSet(100, 200, 300), 300},
	} {
		s := x.s
		max := s.Max()
		if max != x.max {
			t.Errorf("%v.Max() = %d; want %d", s, max, x.max)
		}
	}

	s := NewBitSet()
	if !PanicsByBit((*BitSet).Max, s) {
		t.Errorf("Max() should panic for empty set.")
	}
	CheckInvariants(t, "Max() panic", s)
}

func TestSize(t *testing.T) {
	for _, x := range []struct {
		s    *BitSet
		size int
	}{
		{NewBitSet(), 0},
		{NewBitSet(-1), 0},
		{NewBitSet(1), 1},
		{NewBitSet(64), 1},
		{NewBitSet(65), 1},
		{NewBitSet(1, 2, 3), 3},
		{NewBitSet(100, 200, 300), 3},
		{NewBitSet().AddRange(0, 64), 64},
		{NewBitSet().AddRange(1, 64), 63},
		{NewBitSet().AddRange(0, 63), 63},
	} {
		s := x.s
		size := s.Size()
		if size != x.size {
			t.Errorf("%v.Size() = %d; want %d", s, size, x.size)
		}
	}
}

func TestEmpty(t *testing.T) {
	for _, x := range []struct {
		s     *BitSet
		empty bool
	}{
		{NewBitSet(), true},
		{NewBitSet(-1), true},
		{NewBitSet().AddRange(-10, 0), true},
		{NewBitSet(1), false},
		{NewBitSet(65), false},
		{NewBitSet(1, 2, 3), false},
		{NewBitSet(100, 200, 300), false},
	} {
		s := x.s
		empty := s.Empty()
		if empty != x.empty {
			t.Errorf("%v.Empty() = %v; want %v", s, empty, x.empty)
		}
	}
}

func TestNextPrev(t *testing.T) {
	for _, x := range []struct {
		s     *BitSet
		m     int
		nextN int
		prevN int
	}{
		{NewBitSet(), 1, -1, -1},
		{NewBitSet(), 0, -1, -1},
		{NewBitSet(), -1, -1, -1},

		{NewBitSet(1), -1, 1, -1},
		{NewBitSet(1), 0, 1, -1},
		{NewBitSet(1), 1, -1, -1},
		{NewBitSet(1), 2, -1, 1},

		{NewBitSet(0, 2), -1, 0, -1},
		{NewBitSet(0, 2), 0, 2, -1},
		{NewBitSet(0, 2), 1, 2, 0},
		{NewBitSet(0, 2), 2, -1, 0},
		{NewBitSet(0, 2), 3, -1, 2},

		{NewBitSet(63, 64), 62, 63, -1},
		{NewBitSet(63, 64), 63, 64, -1},
		{NewBitSet(63, 64), 64, -1, 63},
		{NewBitSet(63, 64), 65, -1, 64},

		{NewBitSet(100, 300), MinInt, 100, -1},
		{NewBitSet(100, 300), -1, 100, -1},
		{NewBitSet(100, 300), 0, 100, -1},
		{NewBitSet(100, 300), 1, 100, -1},
		{NewBitSet(100, 300), 99, 100, -1},
		{NewBitSet(100, 300), 100, 300, -1},
		{NewBitSet(100, 300), 101, 300, 100},
		{NewBitSet(100, 300), 200, 300, 100},
		{NewBitSet(100, 300), 299, 300, 100},
		{NewBitSet(100, 300), 300, -1, 100},
		{NewBitSet(100, 300), 301, -1, 300},
		{NewBitSet(100, 300), 400, -1, 300},
		{NewBitSet(100, 300), MaxInt, -1, 300},
	} {
		s := x.s
		m := x.m
		nextN := s.Next(m)
		if nextN != x.nextN {
			t.Errorf("%v.Next(%d) = %d; want %d", s, m, nextN, x.nextN)
		}
		prevN := s.Prev(m)
		if prevN != x.prevN {
			t.Errorf("%v.Prev(%d) = %d; want %d", s, m, prevN, x.prevN)
		}
	}
}

func TestVisit(t *testing.T) {
	for _, x := range []struct {
		s   *BitSet
		res string
	}{
		{NewBitSet(), ""},
		{NewBitSet(0), "0"},
		{NewBitSet(1, 2, 3, 62, 63, 64), "123626364"},
		{NewBitSet(1, 22, 333, 4444), "1223334444"},
	} {
		s := x.s
		res := ""

		s.Visit(func(n int) (skip bool) {
			res += strconv.Itoa(n)
			return
		})
		if res != x.res {
			t.Errorf("%v.Visit(func(n int) { s += Itoa(n) }) -> s=%q; want %q", s, res, x.res)
		}

		s = x.s
		res = ""
		s.Visit(func(n int) (skip bool) {
			s.DeleteRange(0, n+1)
			res += strconv.Itoa(n)
			return
		})
		if res != x.res {
			t.Errorf("%v.Visit(func(n int) { s.DeleteRange(0, n+1); s += Itoa(n) }) -> s=%q; want %q", s, res, x.res)
		}
	}
	s := NewBitSet(1, 2)
	count := 0
	aborted := s.Visit(func(n int) (skip bool) {
		count++
		if n == 1 {
			skip = true
			return
		}
		return
	})
	if aborted == false {
		t.Errorf("Visit returned false when aborted.")
	}
	if count > 1 {
		t.Errorf("Visit didn't abort.")
	}
	count = 0
	aborted = s.Visit(func(n int) (skip bool) {
		count++
		return
	})
	if aborted == true {
		t.Errorf("Visit returned true when not aborted.")
	}
	if count != 2 {
		t.Errorf("Visit aborted.")
	}
}

func TestByString(t *testing.T) {
	for _, x := range []struct {
		s   *BitSet
		res string
	}{
		{NewBitSet(), "{}"},
		{NewBitSet(-1), "{}"},
		{NewBitSet(1), "{1}"},
		{NewBitSet(1, -1), "{1}"},
		{NewBitSet(1, 2), "{1 2}"},
		{NewBitSet(1, 3), "{1 3}"},
		{NewBitSet(0, 2, 3), "{0 2 3}"},
		{NewBitSet(0, 1, 3), "{0 1 3}"},
		{NewBitSet(0, 2, 3, 5), "{0 2 3 5}"},
		{NewBitSet(0, 1, 2, 4, 5), "{0..2 4 5}"},
		{NewBitSet(0, 1, 2, 3, 5, 7, 8, 9), "{0..3 5 7..9}"},
		{NewBitSet(65), "{65}"},
		{NewBitSet(100, 200, 300), "{100 200 300}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("S.String() = %q; want %q", res, x.res)
		}
	}
}

func TestByAdd(t *testing.T) {
	for _, x := range []struct {
		s   *BitSet
		res string
	}{
		{NewBitSet().Add(-1), "{}"},
		{NewBitSet().Add(1), "{1}"},
		{NewBitSet(1).Add(1), "{1}"},
		{NewBitSet(1).Add(2), "{1 2}"},
		{NewBitSet().Add(65), "{65}"},
		{NewBitSet().Add(100).Add(200).Add(300), "{100 200 300}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("s.Add() = %q; want %q", res, x.res)
		}
		CheckInvariants(t, "Add", x.s)
	}
}

func TestDelete(t *testing.T) {
	for _, x := range []struct {
		s   *BitSet
		res string
	}{
		{NewBitSet(1).Delete(1), "{}"},
		{NewBitSet(1).Delete(-1), "{1}"},
		{NewBitSet(1).Delete(2), "{1}"},
		{NewBitSet(65).Delete(64), "{65}"},
		{NewBitSet(100, 200, 300).Delete(200), "{100 300}"},
		{NewBitSet(100, 200, 300).Delete(300), "{100 200}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("s.Delete() = %q; want %q", res, x.res)
		}
		CheckInvariants(t, "Delete", x.s)
	}
}

type rangeFunc struct {
	fInt   func(S *BitSet, n int) *BitSet
	fRange func(S *BitSet, m, n int) *BitSet
	name   string
}

func TestRange(t *testing.T) {
	rangeFuncs := []rangeFunc{
		{(*BitSet).Add, (*BitSet).AddRange, "AddRange"},
		{(*BitSet).Delete, (*BitSet).DeleteRange, "DeleteRange"},
	}

	for _, x := range []struct {
		s    *BitSet
		m, n int
	}{
		{NewBitSet(), 0, 0},
		{NewBitSet(), 2, 1},
		{NewBitSet(), -2, -1},
		{NewBitSet(), -1, 0},
		{NewBitSet(), -1, -1},
		{NewBitSet(), 1, 10},
		{NewBitSet(), 64, 66},
		{NewBitSet(), 1, 1000},

		{NewBitSet(1, 2, 3), 0, 1},
		{NewBitSet(1, 2, 3), 0, 2},
		{NewBitSet(1, 2, 3), 0, 3},
		{NewBitSet(1, 2, 3), 0, 4},
		{NewBitSet(1, 2, 3), 1, 2},
		{NewBitSet(1, 2, 3), 1, 4},
		{NewBitSet(1, 2, 3), 1, 5},
		{NewBitSet(1, 2, 3), 1, 1000},

		{NewBitSet(100, 200, 300), 50, 100},
		{NewBitSet(100, 200, 300), 50, 101},
		{NewBitSet(100, 200, 300), 50, 250},
		{NewBitSet(100, 200, 300), 50, 350},
		{NewBitSet(100, 200, 300), 250, 350},
		{NewBitSet(100, 200, 300), 300, 350},
		{NewBitSet(100, 200, 300), 350, 400},
		{NewBitSet(100, 200, 300), 1, 1000},
	} {
		for _, o := range rangeFuncs {
			fInt, fRange, name := o.fInt, o.fRange, o.name
			s := x.s
			m, n := x.m, x.n

			res := fRange(new(BitSet).Set(s), m, n)
			exp := new(BitSet).Set(s)
			for i := m; i < n; i++ {
				fInt(exp, i)
			}
			if !res.Equal(exp) {
				t.Errorf("%v.%v(%d, %d) = %v; want %v", s, name, m, n, res, exp)
			}
			CheckInvariants(t, name, res)
		}
	}
}

func TestSet(t *testing.T) {
	for _, x := range []struct {
		s, a *BitSet
	}{
		{NewBitSet(), NewBitSet()},
		{NewBitSet(), NewBitSet(1)},
		{NewBitSet(), NewBitSet(65)},
		{NewBitSet(), NewBitSet(1, 2, 3)},
		{NewBitSet(), NewBitSet(100, 200, 300)},

		{NewBitSet(1, 2, 3), NewBitSet()},
		{NewBitSet(1, 2, 3), NewBitSet(1)},
		{NewBitSet(1, 2, 3), NewBitSet(65)},
		{NewBitSet(1, 2, 3), NewBitSet(1, 2, 3)},
		{NewBitSet(1, 2, 3), NewBitSet(100, 200, 300)},

		{NewBitSet(100, 200, 300), NewBitSet()},
		{NewBitSet(100, 200, 300), NewBitSet(1)},
		{NewBitSet(100, 200, 300), NewBitSet(65)},
		{NewBitSet(100, 200, 300), NewBitSet(1, 2, 3)},
		{NewBitSet(100, 300, 300), NewBitSet(100, 200, 300)},
	} {
		s := x.s

		ss := s.Set(x.a)
		if ss != s {
			t.Errorf("&(s.set(%v)) = %p, &S = %p; want same", x.a, ss, s)
		}
		if !ss.Equal(x.a) {
			t.Errorf("s.set(%v) = %v; want %v", x.a, ss, x.a)
		}
		CheckInvariants(t, "set", ss)
	}
}

type binOp struct {
	f    func(s *BitSet, a, b *BitSet) *BitSet
	name string
}

func TestBinOp(t *testing.T) {
	And := binOp{(*BitSet).SetAnd, "SetAnd"}
	AndNot := binOp{(*BitSet).SetAndNot, "SetAndNot"}
	Or := binOp{(*BitSet).SetOr, "SetOr"}
	Xor := binOp{(*BitSet).SetXor, "SetXor"}
	for _, x := range []struct {
		op   binOp
		a, b *BitSet
		exp  *BitSet
	}{
		{And, NewBitSet(), NewBitSet(), NewBitSet()},
		{And, NewBitSet(1), NewBitSet(), NewBitSet()},
		{And, NewBitSet(), NewBitSet(1), NewBitSet()},
		{And, NewBitSet(1), NewBitSet(1), NewBitSet(1)},
		{And, NewBitSet(1), NewBitSet(2), NewBitSet()},
		{And, NewBitSet(1), NewBitSet(1, 2), NewBitSet(1)},
		{And, NewBitSet(1, 2), NewBitSet(2, 3), NewBitSet(2)},
		{And, NewBitSet(100), NewBitSet(), NewBitSet()},
		{And, NewBitSet(), NewBitSet(100), NewBitSet()},
		{And, NewBitSet(100), NewBitSet(100), NewBitSet(100)},
		{And, NewBitSet(100), NewBitSet(100, 200), NewBitSet(100)},
		{And, NewBitSet(200), NewBitSet(100, 200), NewBitSet(200)},
		{And, NewBitSet(100, 200), NewBitSet(200, 300), NewBitSet(200)},

		{AndNot, NewBitSet(), NewBitSet(), NewBitSet()},
		{AndNot, NewBitSet(1), NewBitSet(), NewBitSet(1)},
		{AndNot, NewBitSet(), NewBitSet(1), NewBitSet()},
		{AndNot, NewBitSet(1), NewBitSet(1), NewBitSet()},
		{AndNot, NewBitSet(1), NewBitSet(2), NewBitSet(1)},
		{AndNot, NewBitSet(1), NewBitSet(1, 2), NewBitSet()},
		{AndNot, NewBitSet(1, 2), NewBitSet(2, 3), NewBitSet(1)},
		{AndNot, NewBitSet(100), NewBitSet(), NewBitSet(100)},
		{AndNot, NewBitSet(), NewBitSet(100), NewBitSet()},
		{AndNot, NewBitSet(100), NewBitSet(100), NewBitSet()},
		{AndNot, NewBitSet(100), NewBitSet(100, 200), NewBitSet()},
		{AndNot, NewBitSet(200), NewBitSet(100, 200), NewBitSet()},
		{AndNot, NewBitSet(100, 200), NewBitSet(200, 300), NewBitSet(100)},

		{Or, NewBitSet(), NewBitSet(), NewBitSet()},
		{Or, NewBitSet(), NewBitSet(1), NewBitSet(1)},
		{Or, NewBitSet(1), NewBitSet(), NewBitSet(1)},
		{Or, NewBitSet(1), NewBitSet(1), NewBitSet(1)},
		{Or, NewBitSet(1), NewBitSet(2), NewBitSet(1, 2)},
		{Or, NewBitSet(1), NewBitSet(1, 2), NewBitSet(1, 2)},
		{Or, NewBitSet(1, 2), NewBitSet(2, 3), NewBitSet(1, 2, 3)},
		{Or, NewBitSet(100), NewBitSet(), NewBitSet(100)},
		{Or, NewBitSet(), NewBitSet(100), NewBitSet(100)},
		{Or, NewBitSet(100), NewBitSet(100), NewBitSet(100)},
		{Or, NewBitSet(100), NewBitSet(100, 200), NewBitSet(100, 200)},
		{Or, NewBitSet(200), NewBitSet(100, 200), NewBitSet(100, 200)},
		{Or, NewBitSet(100, 200), NewBitSet(200, 300), NewBitSet(100, 200, 300)},

		{Xor, NewBitSet(), NewBitSet(), NewBitSet()},
		{Xor, NewBitSet(), NewBitSet(1), NewBitSet(1)},
		{Xor, NewBitSet(1), NewBitSet(), NewBitSet(1)},
		{Xor, NewBitSet(1), NewBitSet(1), NewBitSet()},
		{Xor, NewBitSet(1), NewBitSet(2), NewBitSet(1, 2)},
		{Xor, NewBitSet(1), NewBitSet(1, 2), NewBitSet(2)},
		{Xor, NewBitSet(1, 2), NewBitSet(2, 3), NewBitSet(1, 3)},
		{Xor, NewBitSet(100), NewBitSet(), NewBitSet(100)},
		{Xor, NewBitSet(), NewBitSet(100), NewBitSet(100)},
		{Xor, NewBitSet(100), NewBitSet(100), NewBitSet()},
		{Xor, NewBitSet(100), NewBitSet(100, 200), NewBitSet(200)},
		{Xor, NewBitSet(200), NewBitSet(100, 200), NewBitSet(100)},
		{Xor, NewBitSet(100, 200), NewBitSet(200, 300), NewBitSet(100, 300)},
	} {
		op, name := x.op.f, x.op.name
		a, b := NewBitSet().Set(x.a), NewBitSet().Set(x.b)
		s := NewBitSet()

		res := op(s, a, b)
		exp := x.exp
		if s != res {
			t.Errorf("&(s.%s(%v, %v)) = %p &s = %p; want same", name, a, b, s, res)
		}
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = a
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = b
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = NewBitSet().AddRange(150, 250)
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)
	}
}

func TestNextPow2(t *testing.T) {
	for _, x := range []struct {
		n, p int
	}{
		{MinInt, 1},
		{-1, 1},
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 4},
		{4, 8},
		{1<<19 - 1, 1 << 19},
		{1 << 19, 1 << 20},
		{MaxInt >> 1, MaxInt>>1 + 1},
		{MaxInt>>1 + 1, MaxInt},
		{MaxInt - 1, MaxInt},
		{MaxInt, MaxInt},
	} {
		n := x.n

		p := nextPow2(n)
		if p != x.p {
			t.Errorf("nextPow2(%#x) = %#x; want %#x", n, p, x.p)
		}
	}
}
