package meta

import "errors"

var EI_KeyNotImplement error = errors.New("interface Key should be implements")
var EI_StringNotImplement error = errors.New("interface String should be implements")
var EI_ValueNotImplement error = errors.New("interface Value should be implements")

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

// String 读取对象的String表现
type String interface {
	String() string
}

// Value 使用字符串初始化该对象
type Value interface {
	Value(v string) error
}
