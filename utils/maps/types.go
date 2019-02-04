package maps

type IMap interface {
	Get(k interface{}) (val interface{}, ok bool)
	Set(k, v interface{})
	Contains(k interface{}) bool
	Keys() []interface{}
	Delete(k interface{})
	ForEach(func(k, v interface{}))
	Clear()
}
