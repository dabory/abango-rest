// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abango

import (
	"context"
	"time"

	"github.com/dabory/abango-rest/etc"
	e "github.com/dabory/abango-rest/etc"
	"github.com/go-redis/redis/v8"
	// "github.com/tidwall/buntdb"
)

var (
	// ADB *buntdb.DB //AegisCache

	RedisCtx = context.Background()
	MDB      *redis.Client //QDB은 Redis를 도입하면서 통합함.
	// QDB *buntdb.DB
	QDBOn bool // QDb에서 쿼리 가져옴
)

func MdbView(key string) (retval string, reterr error) {

	value, err := MDB.Get(RedisCtx, key).Result()
	if err == redis.Nil {
		reterr = e.LogErr("ASDF1QWERCAA", "MDB.View Not Found in Key: "+key, err)
	} else if err != nil {
		reterr = e.LogErr("ASDFQWERA", "MDB.View Error reading data: "+key, err)
	}
	return value, reterr
}

func MdbUpdate(key string, value string) (reterr error) {

	REDIS_EXTIME := 12 * time.Hour
	err := MDB.Set(RedisCtx, key, value, REDIS_EXTIME).Err()
	if err != nil {
		reterr = e.MyErr("QWVGAVAEFV-MDB.Update Error in Key: "+key+" Value: "+value, err, false)
	}
	return nil
}

func MdbDelete(key string) error {
	err := MDB.Del(RedisCtx, key).Err()
	if err != nil {
		return e.MyErr("QWVGAVAEFV-MDB.Delete Error in Key: "+key, err, false)
	}
	return nil
}

// func QdbView(key string) (retval string, reterr error) {

// 	QDB.View(func(tx *buntdb.Tx) error {
// 		if value, err := tx.Get(key); err == nil {
// 			retval = value
// 			reterr = nil
// 		} else {
// 			retval = ""
// 			reterr = errors.New("QDB.View Not Found in Key: " + key)
// 		}
// 		return nil
// 	})
// 	return retval, reterr
// }

// func QdbUpdate(key string, value string) (reterr error) {

// 	QDB.Update(func(tx *buntdb.Tx) error {
// 		_, _, err := tx.Set(key, value, nil)
// 		if err != nil {
// 			reterr = e.MyErr("TKBKUYIH-QDB.Update Error in Key: "+key+" Value: "+value, err, false)
// 		}
// 		return nil
// 	})
// 	return nil
// }

func GetQryStr(filename string) (string, error) {
	var str string
	var err error

	if QDBOn {
		if str, err = MdbView(filename); err == nil {
			return str, nil
		}
	}

	// 공통 경로: 파일에서 로딩
	if str, err = e.FileToQryChkStr(filename); err != nil {
		return "", etc.LogErr("PKOJHKJUY", "File", err)
	}

	// QDBOn인 경우에만 메모리에 저장
	if QDBOn {
		if err := MdbUpdate(filename, str); err != nil {
			return "", etc.LogErr("OIUJLJOUJLH", "QdbUpdate Failed", err)
		}
	}

	return str, nil
}

// func GetQryStr(filename string) (string, error) {

// 	var str string
// 	var err error
// 	if QDBOn {
// 		if str, err = MdbView(filename); err == nil {
// 			// etc.LogNil("Qry from Memory!!")
// 			return str, nil
// 		} else {
// 			if str, err = e.FileToQryChkStr(filename); err == nil {
// 				if err := MdbUpdate(filename, str); err != nil {
// 					return "", etc.LogErr("OIUJLJOUJLH", "QdbUpdate Failed ", err)
// 				}
// 				// etc.LogNil("Qry from File!!")
// 				return str, nil
// 			} else {
// 				return "", err
// 			}
// 		}
// 	} else {
// 		if str, err = e.FileToQryChkStr(filename); err == nil {
// 			// etc.LogNil("QRY FILE")
// 			return str, nil
// 		} else {
// 			return "", etc.LogErr("PKOJHKJUY", " File", err)
// 		}
// 	}
// }
