package storage

// Iterator 迭代函数
type Iterator func(interface{}) error

// ForEach 对象遍历
type ForEach interface {
	ForEach(Iterator) error
}
