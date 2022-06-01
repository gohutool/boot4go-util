package data

import (
	"errors"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : loop_queue.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/31 21:27
* 修改历史 : 1. [2022/5/31 21:27] 创建文件 by LongYong
*/

var (
	// errQueueIsFull will be returned when the worker queue is full.
	errQueueIsFull = errors.New("the queue is full")

	// errQueueIsReleased will be returned when trying to insert item to a released worker queue.
	errQueueIsReleased = errors.New("the queue length is zero")
)

type CacheUnit[T any] struct {
	obj              *T
	lastUsedDuration time.Time
}

type loopQueue[T any] struct {
	items  []*CacheUnit[T]
	expiry []*CacheUnit[T]
	head   int
	tail   int
	size   int
	isFull bool
}

func NewLoopQueue[T any](size int) *loopQueue[T] {
	return &loopQueue[T]{
		items: make([]*CacheUnit[T], size),
		size:  size,
	}
}

func (lq *loopQueue[T]) Len() int {
	if lq.size == 0 {
		return 0
	}

	if lq.head == lq.tail {
		if lq.isFull {
			return lq.size
		}
		return 0
	}

	if lq.tail > lq.head {
		return lq.tail - lq.head
	}

	return lq.size - lq.head + lq.tail
}

func (lq *loopQueue[T]) IsEmpty() bool {
	return lq.head == lq.tail && !lq.isFull
}

func (lq *loopQueue[T]) Insert(t *T) error {
	if lq.size == 0 {
		return errQueueIsReleased
	}

	if lq.isFull {
		return errQueueIsFull
	}

	if lq.items[lq.tail] == nil {
		lq.items[lq.tail] = &CacheUnit[T]{obj: t}
	} else {
		lq.items[lq.tail].obj = t
	}

	lq.items[lq.tail].lastUsedDuration = time.Now()

	lq.tail++

	if lq.tail == lq.size {
		lq.tail = 0
	}
	if lq.tail == lq.head {
		lq.isFull = true
	}

	return nil
}

func (lq *loopQueue[T]) Get() *T {
	if lq.IsEmpty() {
		return nil
	}

	w := lq.items[lq.head]

	return w.obj
}

func (lq *loopQueue[T]) Detach() *T {
	if lq.IsEmpty() {
		return nil
	}

	w := lq.items[lq.head]
	lq.items[lq.head] = nil
	lq.head++
	if lq.head == lq.size {
		lq.head = 0
	}
	lq.isFull = false

	return w.obj
}

func (lq *loopQueue[T]) RetrieveExpiry(duration time.Duration) []*T {
	expiryTime := time.Now().Add(-duration)
	index := lq.binarySearch(expiryTime)
	if index == -1 {
		return nil
	}
	lq.expiry = lq.expiry[:0]

	var v []*CacheUnit[T] = lq.expiry

	if lq.head <= index {
		lq.expiry = append(v, lq.items[lq.head:index+1]...)
		for i := lq.head; i < index+1; i++ {
			lq.items[i] = nil
		}
	} else {
		lq.expiry = append(v, lq.items[0:index+1]...)
		lq.expiry = append(v, lq.items[lq.head:]...)
		for i := 0; i < index+1; i++ {
			lq.items[i] = nil
		}
		for i := lq.head; i < lq.size; i++ {
			lq.items[i] = nil
		}
	}
	head := (index + 1) % lq.size
	lq.head = head
	v = lq.expiry
	if len(v) > 0 {
		lq.isFull = false
	}

	var rtn []*T

	for _, o := range lq.expiry {
		rtn = append(rtn, o.obj)
	}
	return rtn
}

func (lq *loopQueue[T]) binarySearch(expiryTime time.Time) int {
	var o []*CacheUnit[T] = lq.items
	var mid, nlen, basel, tmid int
	nlen = len(o)

	// if no need to remove work, return -1
	if lq.IsEmpty() || expiryTime.Before(lq.items[lq.head].lastUsedDuration) {
		return -1
	}

	// example
	// size = 8, head = 7, tail = 4
	// [ 2, 3, 4, 5, nil, nil, nil,  1]  true position
	//   0  1  2  3    4   5     6   7
	//              tail          head
	//
	//   1  2  3  4  nil nil   nil   0   mapped position
	//            r                  l

	// base algorithm is a copy from worker_stack
	// map head and tail to effective left and right
	r := (lq.tail - 1 - lq.head + nlen) % nlen
	basel = lq.head
	l := 0
	for l <= r {
		mid = l + ((r - l) >> 1)
		// calculate true mid position from mapped mid position
		tmid = (mid + basel + nlen) % nlen
		if expiryTime.Before(lq.items[tmid].lastUsedDuration) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	// return true position from mapped position
	return (r + basel + nlen) % nlen
}

func (lq *loopQueue[T]) reset() {
	if lq.IsEmpty() {
		return
	}

	for i := 0; i < lq.Len(); i++ {
		lq.items[i].obj = nil
		lq.items[i] = nil
	}
	lq.items = lq.items[:0]
	lq.size = 0
	lq.head = 0
	lq.tail = 0

	//Releasing:
	//	if w := lq.detach(); w != nil {
	//		w = nil
	//		goto Releasing
	//	}
	//	lq.items = lq.items[:0]
	//	lq.size = 0
	//	lq.head = 0
	//	lq.tail = 0
}
