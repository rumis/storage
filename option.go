package storage

// ForEach 对象遍历
type ForEach interface {
	ForEach(fn func(interface{}) error) error
}
