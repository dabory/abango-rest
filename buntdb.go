// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abango

import (
	e "github.com/dabory/abango-rest/etc"
	"github.com/tidwall/buntdb"
)

func MdbView(key string) (retval string, reterr error) {

	MDB.View(func(tx *buntdb.Tx) error {
		if value, err := tx.Get(key); err == nil {
			retval = value
			reterr = nil
		} else {
			retval = ""
			reterr = e.LogErr("ASDFQWERCAA", "MDB.View Not Fount in Key: "+key, err)
		}
		return nil
	})
	return retval, reterr
}

func MdbUpdate(key string, value string) (reterr error) {

	MDB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		if err != nil {
			reterr = e.MyErr("QWVGAVAEFV-MDB.Update Error in Key: "+key+" Value: "+value, err, false)
		}
		return nil
	})
	return nil
}
