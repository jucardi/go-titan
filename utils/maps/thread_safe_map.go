package maps

import (
	"sync"
)

type threadSafeMap struct {
	m  map[interface{}]interface{}
	mx sync.RWMutex
}

func NewThreadSafeMap() IMap {
	return &threadSafeMap{
		m:  make(map[interface{}]interface{}),
		mx: sync.RWMutex{},
	}
}

func (r *threadSafeMap) Get(k interface{}) (val interface{}, ok bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	val, ok = r.m[k]
	return
}

func (r *threadSafeMap) Set(k, v interface{}) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.m[k] = v
}

func (r *threadSafeMap) Contains(k interface{}) bool {
	_, ok := r.Get(k)
	return ok
}

func (r *threadSafeMap) Keys() []interface{} {
	r.mx.RLock()
	defer r.mx.RUnlock()
	var ret []interface{}
	for k := range r.m {
		ret = append(ret, k)
	}
	return ret
}

func (r *threadSafeMap) Delete(k interface{}) {
	r.mx.Lock()
	defer r.mx.Unlock()
	delete(r.m, k)
}

func (r *threadSafeMap) ForEach(fn func(k, v interface{})) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	for k, v := range r.m {
		fn(k, v)
	}
}

func (r *threadSafeMap) Clear() {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.m = make(map[interface{}]interface{})
}
