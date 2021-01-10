package abango

import (
	_ "github.com/go-sql-driver/mysql"
)

func RestSvcStandBy(RouterHandler func(*AbangoAsk)) {

	// MyLinkXDB()
	var v AbangoAsk
	RouterHandler(&v)
}
