// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"reflect"
	"sync"
)

type pool struct {
	size int
	tp   reflect.Type
	stp  reflect.Type
	pool reflect.Value
	cur  int
	lock sync.Mutex
}

func newPool(size int, tp reflect.Type) *pool {
	return &pool{
		size: size,
		cur:  0,
		pool: reflect.MakeSlice(reflect.SliceOf(tp), 0, 0), // init for don't allocate memory
		tp:   reflect.SliceOf(tp),
		stp:  tp,
	}
}

func (p *pool) New() reflect.Value {
	//return reflect.New(p.stp)
	p.lock.Lock()
	if p.cur == p.pool.Len() {
		p.pool = reflect.MakeSlice(p.tp, p.size, p.size)
		p.cur = 0
	}
	res := p.pool.Index(p.cur).Addr()
	p.cur++
	p.lock.Unlock()
	return res
}
