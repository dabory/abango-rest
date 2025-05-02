package abgqry

import (
	"errors"

	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
)

func OneRowQry(y *abango.Controller, sql string) (c1 string, c2 string, c3 string, c4 string, c5 string, err error) {
	page, err := y.Db.Query(sql)
	if err != nil {
		return "", "", "", "", "", e.LogErr("so9eowad3", e.FuncNameErr()+":"+sql+"\n", err)
	}
	if len(page) > 1 {
		return "", "", "", "", "", e.LogErr("so982wds3", e.FuncNameErr(), errors.New("Row Count > 1 :"+sql+"\n"))
	}

	for _, row := range page {
		c1 = string(row["c1"])
		c2 = string(row["c2"])
		c3 = string(row["c3"])
		c4 = string(row["c4"])
		c5 = string(row["c5"])
	}
	return
}

func OneRowQry10(y *abango.Controller, sql string) (c1 string, c2 string, c3 string,
	c4 string, c5 string, c6 string, c7 string, c8 string, c9 string, c10 string, err error) {
	page, err := y.Db.Query(sql)
	if err != nil {
		return "", "", "", "", "", "", "", "", "", "", e.LogErr("so9eowd3", e.FuncNameErr()+":"+sql+"\n", err)
	}
	if len(page) > 1 {
		return "", "", "", "", "", "", "", "", "", "", e.LogErr("so982wd3", e.FuncNameErr(), errors.New("Row Count > 1 :"+sql+"\n"))
	}

	for _, row := range page {
		c1 = string(row["c1"])
		c2 = string(row["c2"])
		c3 = string(row["c3"])
		c4 = string(row["c4"])
		c5 = string(row["c5"])
		c6 = string(row["c6"])
		c7 = string(row["c7"])
		c8 = string(row["c8"])
		c9 = string(row["c9"])
		c10 = string(row["c10"])
	}
	return
}
