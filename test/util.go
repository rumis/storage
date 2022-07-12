package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rumis/seal"
	"github.com/rumis/seal/builder"
	"github.com/rumis/storage/scache"
	"github.com/rumis/storage/srepo"
)

// InitClient 初始化
func InitClient() sqlmock.Sqlmock {
	// 启动内存Redis服务并创建Client
	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})
	scache.SetDefaultClient(rClient)

	// 数据库
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(err)
	}
	sealDb, err := seal.OpenWithDB(db, builder.NewMysqlBuilder())
	if err != nil {
		panic(err)
	}
	srepo.SetSealR(sealDb)
	srepo.SetSealW(sealDb)

	return mock
}
