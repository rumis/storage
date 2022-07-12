package meta

type OptionStatus int8

const (
	OptionStatusContinue OptionStatus = 1 // 流程未完成，需要继续执行
	OptionStatusBreak    OptionStatus = 2 // 流程已完成，可以返回结果了
)
