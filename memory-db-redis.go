// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abango

import (
	"context"
	"time"

	e "github.com/dabory/abango-rest/etc"
	"github.com/go-redis/redis/v8"
	// "github.com/tidwall/buntdb"
)

var (
	// ADB *buntdb.DB //AegisCache

	RedisCtx = context.Background()
	QDB      *redis.Client //QDB은 Redis를 도입하면서 통합함.
	// QDB *buntdb.DB
	QDBOn bool // QDb에서 쿼리 가져옴
)

func QdbView(key string) (retval string, reterr error) {

	value, err := QDB.Get(RedisCtx, key).Result()
	if err == redis.Nil {
		reterr = e.LogErr("ASDF1QWERCAA", "QDB.View Not Found in Key: "+key, err)
	} else if err != nil {
		reterr = e.LogErr("ASDFQWERA", "QDB.View Error reading data: "+key, err)
	}
	return value, reterr
}

func QdbUpdate(key string, value string) (reterr error) {

	REDIS_EXTIME := 12 * time.Hour
	err := QDB.Set(RedisCtx, key, value, REDIS_EXTIME).Err()
	if err != nil {
		reterr = e.MyErr("QWVGAVAEFV-QDB.Update Error in Key: "+key+" Value: "+value, err, false)
	}
	return nil
}

func QdbDelete(key string) error {
	err := QDB.Del(RedisCtx, key).Err()
	if err != nil {
		return e.MyErr("QWVGAVAEFV-QDB.Delete Error in Key: "+key, err, false)
	}
	return nil
}
