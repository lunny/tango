package tango

import (
	"sync"
	"reflect"
)

type pools struct {
	pools map[reflect.Type]*pool
	lock sync.Mutex
	size int
}

func NewPools(size int) *pools {
	return &pools{
		pools: make(map[reflect.Type]*pool),
		size: size,
	}
}

func (ps *pools) Pool(tp reflect.Type) *pool {
	var p *pool
	var ok bool
	ps.lock.Lock()
	if p, ok = ps.pools[tp]; !ok {
		p = newPool(ps.size, tp)
		ps.pools[tp] = p
	}
	ps.lock.Unlock()
	return p
}

func (p *pools) New(tp reflect.Type) reflect.Value {
	return p.Pool(tp).New()
}

type pool struct {
	size int
	tp reflect.Type
	pool reflect.Value
	cur int
	lock sync.Mutex
}

func newPool(size int, tp reflect.Type) *pool {
	return &pool{
		size: size,
		cur: 0,
		pool: reflect.MakeSlice(reflect.SliceOf(tp), size, size),
		tp: reflect.SliceOf(tp),
	}
}

func (p *pool) New() reflect.Value {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.cur == p.pool.Len() {
		p.pool = reflect.MakeSlice(p.tp, p.size, p.size)
		p.cur = 0
	}
	p.cur++
	return p.pool.Index(p.cur).Addr()
}