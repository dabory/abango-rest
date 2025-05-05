package abg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
	"github.com/go-xorm/xorm"
)

func QryCount(YDB *xorm.Engine, sql string) int64 { // 이렇게 error 없이 가는 건 아주 특이한 케이스이다.

	arr, err := YDB.Query(sql)
	if err != nil {
		e.ErrLog("model.QryCount Failure !\n["+sql+"]\n", err)
		return 0
	}
	if len(arr) == 1 {
		for _, buf := range arr[0] {
			cnt, _ := strconv.Atoi(string(buf))
			return int64(cnt)
		}
	} else {
		e.LogErr("0237ew32", e.FuncNameErr(), errors.New("QryCount Line is more than 1"))
		return 0
	}
	return 0
}

func SlipIdGet(y *abango.Controller, table, slipPrefix, slipno string) (int, error) {
	if table == "" || slipno == "" {
		return 0, e.LogErr("0237hld92", e.FuncNameErr(), fmt.Errorf("table or slipno is EMPTY"))
	}

	columnName := fmt.Sprintf("%s_no", slipPrefix)
	query := fmt.Sprintf("SELECT id AS c1 FROM %s WHERE %s = '%s'", table, columnName, slipno)

	id, _, _, _, _, _ := OneRowQry(y, query)
	if id == "" {
		return 0, e.LogErr("mclr9uol9", e.FuncNameErr(), fmt.Errorf("Query Error:\n%s", query))
	}

	return e.StrToInt(id), nil
}

func QryDirName(y *abango.Controller, qryName string) (string, string) {

	proDir := ""
	qpName := ""
	if strings.Contains(qryName, "pro:") { // 나중에 이거 Deprecate 시킬 것.
		proDir = "/pro"
		qpName = strings.Replace(qryName, "pro:", "", 1)
	} else if strings.Contains(qryName, "myapp:") {
		proDir = "/myapp"
		qpName = strings.Replace(qryName, "myapp:", "", 1)
	} else {
		qpName = qryName
		proDir = "/erp"
	}

	if !strings.Contains(qryName, "::") {
		return QHOME_DIR + proDir, qpName
	} else {
		q := strings.Split(qpName, "::")
		return strings.ReplaceAll(THEME_QRY_DIR, "$$theme", q[0]) + proDir, q[1]
	}
}

func LastQry(qry xorm.Session) string {
	ret, _ := qry.LastSQL()
	fmt.Println("\n" + ret + "\n")
	return ret
}

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

func SumToBal(y *abango.Controller, qName string, bdId int) error {

	fmt.Println(e.FuncNameInfo(), qName)
	sqlFile := QHOME_DIR + "/erp" + COPY_DIR + "/sum-to-bal/" + qName + ".sql" // theme query 허용하지 않음.
	sqlExec, err := abango.GetQryStr(sqlFile)
	if err != nil {
		return err
	}
	sqlExec = fmt.Sprintf(sqlExec, bdId, bdId)
	if _, err := y.Db.Exec(sqlExec); err != nil {
		return e.LogErr("lskj34r39jsy", e.FuncNameErr()+e.QryPathSql(sqlFile, sqlExec), err)
	}
	return nil
}

func OnDelRestrict(y *abango.Controller, myIdNm, relHdIdNm, relTbl, hdNoNm, hdTbl string, myId int) error {
	// fmt.Println("JJJJ", myIdNm, relHdIdNm, relTbl, hdNoNm, hdTbl, myId)
	// 후속 레코드 여부 판별
	myStr := fmt.Sprintf(`
			SELECT 
				IFNULL(COUNT(*), 0) AS c1,
				IFNULL(MIN(%s), '') AS c2
			FROM %s
			WHERE %s = %d;
	`, relHdIdNm, relTbl, myIdNm, myId)
	// myStr := fmt.Sprintf(`
	// 	SELECT
	// 		COALESCE(c1, 0) AS c1,
	// 		COALESCE(c2, '') AS c2
	// 	FROM (
	// 		SELECT
	// 			COUNT(*) AS c1,
	// 			MIN(%s) AS c2
	// 		FROM %s
	// 		WHERE %s = %d
	// 	) AS grouped_data
	// 	RIGHT JOIN (SELECT 1 AS dummy) AS fallback ON 1 = 1;
	// `, relHdIdNm, relTbl, myIdNm, myId)

	myPreCnt, hdIdStr, _, _, _, err := OneRowQry(y, myStr)
	if err != nil {
		return e.LogErr("sd2asdf", e.FuncNameErr()+myStr, err)
	}

	// 에러메시지 표현 구성 - 사용자가 메시지를 보고 삭제할 전표를 찾을 수 있어야 한다.
	if myPreCnt != "0" {
		hdStr := fmt.Sprintf(`
			SELECT 
				IFNULL(%s, '') AS c1
			FROM %s
			WHERE id = %d
			LIMIT 1;
		`, hdNoNm, hdTbl, e.StrToInt(hdIdStr)) // 위의 LIMIT 1 은 정확히 1 Row 보장.
		// hdStr := fmt.Sprintf(`
		// 	SELECT
		// 		COALESCE(c1, '') AS c1
		// 	FROM (
		// 		SELECT
		// 			%s AS c1
		// 		FROM %s
		// 		WHERE id = %d
		// 	) AS selected_data
		// 	RIGHT JOIN (SELECT 1 AS dummy) AS fallback ON 1 = 1;
		// `, hdNoNm, hdTbl, e.StrToInt(hdIdStr))
		hdNo, _, _, _, _, err := OneRowQry(y, hdStr)
		if err != nil {
			return e.LogErr("sd2asdf4", e.FuncNameErr()+myStr, err)
		} else {
			return fmt.Errorf("OnDeleteRestrict: " + hdNo + " Exists in " + e.TransTblName(hdTbl))
		}
	} else {
		return nil
	}
}
