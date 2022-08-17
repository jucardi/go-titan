package mocking

import "reflect"

type IMock interface {
	When(method string, condition func(args ...interface{}) bool, rets []interface{})
	WhenDelayedRets(method string, condition func(args ...interface{}) bool, retProvider func() []interface{})
	Times(method string) int
	Invoke(method string, args ...interface{}) (ret []interface{})
	ClearWhen()
	ClearTimes()
	Clear()
}

type MockBase struct {
	conditions map[string][]*conditionSet
	times      map[string]int
	returns    map[string][]reflect.Type
}

type conditionSet struct {
	match    func(args ...interface{}) bool
	rets     []interface{}
	provider func() []interface{}
}

func (m *MockBase) Init(interfaceRef interface{}) *MockBase {
	var t reflect.Type
	if x, ok := interfaceRef.(reflect.Type); ok {
		t = x
	} else {
		t = reflect.TypeOf(interfaceRef)
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		var types []reflect.Type
		for j := 0; j < method.Type.NumOut(); j++ {
			types = append(types, method.Type.Out(j))
		}
		m.returns[method.Name] = types
	}
	return m
}

func (m *MockBase) When(method string, condition func(args ...interface{}) bool, rets []interface{}) {
	if _, ok := m.conditions[method]; !ok {
		m.conditions[method] = []*conditionSet{}
	}
	m.conditions[method] = append(m.conditions[method], &conditionSet{
		match: condition,
		rets:  rets,
	})
}

func (m *MockBase) WhenDelayedRets(method string, condition func(args ...interface{}) bool, retProvider func() []interface{}) {
	if _, ok := m.conditions[method]; !ok {
		m.conditions[method] = []*conditionSet{}
	}
	m.conditions[method] = append(m.conditions[method], &conditionSet{
		match:    condition,
		provider: retProvider,
	})
}

func (m *MockBase) Times(method string) int {
	if t, ok := m.times[method]; ok {
		return t
	}
	return 0
}

func (m *MockBase) Clear() {
	m.ClearWhen()
	m.ClearTimes()
}

func (m *MockBase) ClearWhen() {
	m.conditions = map[string][]*conditionSet{}
}

func (m *MockBase) ClearTimes() {
	m.times = map[string]int{}
}

func (m *MockBase) Invoke(method string, args ...interface{}) (ret []interface{}) {
	m.times[method]++
	if set, ok := m.conditions[method]; ok {
		for _, c := range set {
			if c.match(args...) {
				if c.provider != nil {
					return c.provider()
				}
				return c.rets
			}
		}
	}

	if returns, ok := m.returns[method]; ok {
		for _, r := range returns {
			ret = append(ret, reflect.Zero(r).Interface())
		}
	}
	return ret
}

func NewMock() *MockBase {
	return &MockBase{
		conditions: map[string][]*conditionSet{},
		times:      map[string]int{},
		returns:    map[string][]reflect.Type{},
	}
}
