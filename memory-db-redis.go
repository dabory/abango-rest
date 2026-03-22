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
	RedisCtx     = context.Background()
	QDB          *redis.Client //QDB은 Redis를 도입하면서 통합함.
	QDBOn        bool          // QDb에서 쿼리 가져옴
	REDIS_EXTIME = 24 * time.Hour
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

func QdbUpdate(key string, value string) error {
	err := QDB.Set(RedisCtx, key, value, REDIS_EXTIME).Err()
	if err != nil {
		return e.LogErr("ASDFQWERA", e.FuncNameErr()+": QDB.Update Error : "+key, err)
	}
	return nil
}

// func QdbUpdate(key string, value string) (reterr error) {

// 	fmt.Println("key:", key)
// 	fmt.Println("key:", value)

// 	err := QDB.Set(RedisCtx, key, value, REDIS_EXTIME).Err()
// 	if err != nil {
// 		reterr = e.LogErr("ASDFQWERA", e.FuncNameErr()+": QDB.Update Error : "+key, err)
// 	}
// 	return nil
// }

func QdbDelete(key string) error {
	err := QDB.Del(RedisCtx, key).Err()
	if err != nil {
		return e.MyErr("QWVGAVAEFV-QDB.Delete Error in Key: "+key, err, false)
	}
	return nil
}
