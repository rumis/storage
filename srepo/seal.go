package srepo

import (
	"errors"
	"sync"

	"github.com/rumis/seal"
)

var sealR seal.DB
var sealW seal.DB
var sealDBs map[string]seal.DB
var sealDbOnce sync.Once

// SealR 获取Seal只读实例(从库)
func SealR() seal.DB {
	return sealR
}

// SealW 获取Seal读写实例(主库)
func SealW() seal.DB {
	return sealW
}

// SetSealR 设置Seal只读实例(从库)
func SetSealR(db seal.DB) {
	sealR = db
}

// SetSealW 设置Seal读写实例(主库)
func SetSealW(db seal.DB) {
	sealW = db
}

// SetSealDB DB池添加实例
func SetSealDB(inst string, db seal.DB) {
	sealDbOnce.Do(func() {
		sealDBs = make(map[string]seal.DB)
	})
	sealDBs[inst] = db
}

// GetSealDB 获取池中DB实例
func GetSealDB(inst string) (seal.DB, error) {
	db, ok := sealDBs[inst]
	if !ok {
		return seal.DB{}, errors.New("not found")
	}
	return db, nil
}
