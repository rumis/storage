package meta

type DataIntactStat int8

const (
	DataStatCacheAll  DataIntactStat = 1 // 缓存中包含全部数据
	DataStatCachePart DataIntactStat = 2 // 缓存中包含部分数据
	DataStatRepoAll   DataIntactStat = 3 // 缓存+库包含全部数据
	DataStatRepoPart  DataIntactStat = 4 // 缓存+库包含部分数据
)

type OptionStatus int8

const (
	OptionStatusContinue = 1 // 流程未完成，需要继续执行
	OptionStatusBreak    = 2 // 流程已完成，可以返回结果了
)
