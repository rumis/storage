package meta

// Iterator 迭代函数
type Iterator func(interface{}) error

// ForEach 对象遍历
type ForEach interface {
	ForEach(Iterator) error
}

// Zero 判定对象是否为零值
type Zero interface {
	Zero() bool
}

// Key 获取对象的KEY
type Key interface {
	Key() string
}
