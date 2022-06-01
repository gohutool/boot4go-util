package data

import "time"

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : stack.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/31 22:02
* 修改历史 : 1. [2022/5/31 22:02] 创建文件 by LongYong
*/

type stack[T any] struct {
	items  []*CacheUnit[T]
	expiry []*CacheUnit[T]
	size   int
}

func NewStack[T any](size int) *stack[T] {

	return &stack[T]{
		items: make([]*CacheUnit[T], 0, size),
		size:  size,
	}
}

func (s *stack[T]) Len() int {
	var o []*CacheUnit[T] = s.items
	return len(o)
}

func (s *stack[T]) IsEmpty() bool {
	var o []*CacheUnit[T] = s.items
	return len(o) == 0
}

func (s *stack[T]) Insert(t *T) error {
	var o []*CacheUnit[T] = s.items
	var v *CacheUnit[T] = &CacheUnit[T]{obj: t, lastUsedDuration: time.Now()}
	s.items = append(o, v)
	return nil
}

func (s *stack[T]) Get() *T {
	l := s.Len()
	if l == 0 {
		return nil
	}

	w := s.items[l-1]

	return w.obj
}

func (s *stack[T]) Detach() *T {
	l := s.Len()
	if l == 0 {
		return nil
	}

	w := s.items[l-1]
	s.items[l-1] = nil // avoid memory leaks
	s.items = s.items[:l-1]

	return w.obj
}

func (s *stack[T]) RetrieveExpiry(duration time.Duration) []*T {
	n := s.Len()
	if n == 0 {
		return nil
	}

	expiryTime := time.Now().Add(-duration)
	index := s.binarySearch(0, n-1, expiryTime)

	var v []*CacheUnit[T] = s.expiry

	s.expiry = s.expiry[:0]
	if index != -1 {
		s.expiry = append(v, s.items[:index+1]...)

		v = s.items
		m := copy(v, s.items[index+1:])
		for i := m; i < n; i++ {
			s.items[i] = nil
		}
		s.items = s.items[:m]
	}

	var rtn []*T

	for _, o := range s.expiry {
		rtn = append(rtn, o.obj)
	}

	return rtn
}

func (s *stack[T]) binarySearch(l, r int, expiryTime time.Time) int {
	var mid int
	for l <= r {
		mid = (l + r) / 2
		if expiryTime.Before(s.items[mid].lastUsedDuration) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return r
}

func (s *stack[T]) Reset() {
	for i := 0; i < s.Len(); i++ {
		s.items[i].obj = nil
		s.items[i] = nil
	}
	s.items = s.items[:0]
}
