// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abango

import (
	"github.com/dabory/abango-rest/etc"
	e "github.com/dabory/abango-rest/etc"
	"github.com/tidwall/buntdb"
)

var (
	MDB        *buntdb.DB
	QDB        *buntdb.DB
	QRYfromQDB bool
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

func QdbView(key string) (retval string, reterr error) {

	QDB.View(func(tx *buntdb.Tx) error {
		if value, err := tx.Get(key); err == nil {
			retval = value
			reterr = nil
		} else {
			retval = ""
			reterr = e.LogErr("EHGVOIHLHNLM", "QDB.View Not Fount in Key: "+key, err)
		}
		return nil
	})
	return retval, reterr
}

func QdbUpdate(key string, value string) (reterr error) {

	QDB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		if err != nil {
			reterr = e.MyErr("TKBKUYIH-QDB.Update Error in Key: "+key+" Value: "+value, err, false)
		}
		return nil
	})
	return nil
}

func MemoryToQryStr(filename string) (string, error) {

	var str string
	var err error
	if QRYfromQDB {
		if str, err = QdbView(filename); err == nil {
			etc.OkLog("Qry from Memory!!")
			return str, nil
		} else {
			if str, err = e.FileToStrSkip(filename); err == nil {
				if err := QdbUpdate(filename, str); err != nil {
					return "", etc.LogErr("OIUJLJOUJLH", "QdbUpdate Failed ", err)
				}
				etc.OkLog("Qry from File!!")
				return str, nil
			} else {
				return "", etc.LogErr("OKMYFDER", filename+" file does NOT exist.", err)
			}
		}
	} else {
		if str, err = e.FileToStrSkip(filename); err == nil {
			etc.OkLog("QRY FILE")
			return str, nil
		} else {
			return "", etc.LogErr("PKOJHKJUY", filename+" file does NOT exist.", err)
		}
	}

}
