package abg

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
)

func GetCt1WhereCnd(y *abango.Controller, lt *ListType1Vars, sqlStr string) string {
	where := QkWhere

	filterConditions := []struct {
		column string
		start  string
		end    string
	}{
		{lt.FilterDate, lt.StartDate, lt.EndDate},
		{lt.FilterFirst, lt.StartFirst, lt.EndFirst},
		{lt.FilterSecond, lt.StartSecond, lt.EndSecond},
		{lt.FilterThird, lt.StartThird, lt.EndThird},
		{lt.FilterFourth, lt.StartFourth, lt.EndFourth},
	}

	for _, cond := range filterConditions {
		if cond.column != "" && cond.start != "" && cond.end != "" {
			where += fmt.Sprintf(" AND %s BETWEEN '%s' AND '%s' ", cond.column, cond.start, cond.end)
		}
	}

	where += GetExtractStr(y, sqlStr)
	return where
}

// 안쓰니까 Deprecate 처리할 것.
func GetSt1Cnd(y *abango.Controller, st *SelectFilters, sqlStr string) string {

	var where string
	pf := st.Prefix + "."
	for key, str := range st.Str {
		if str.FilterValue != "" {
			where += " and " + pf + "str" + strconv.Itoa(key) + " = '" + str.FilterValue + "' "
		}
	}

	for i, chk := range st.Chk {
		where += " and ("
		for j, opt := range chk.Opt {
			if opt.FilterValue != "" {
				where += " " + pf + "chk" + strconv.Itoa(i) + "_opt" + strconv.Itoa(j) + " = '" + opt.FilterValue + "' or "
			} else {
				where += " true or "
			}
		}
		where += " false ) "
	}

	for key, str := range st.Rng {
		if str.FromValue != "" && str.ToValue != "" {
			where += " and " + pf + "rng" + strconv.Itoa(key) + " between '" + str.FromValue + "' and '" + str.ToValue + "' "
		}
	}

	for key, str := range st.Dec {
		if str.FromValue != "" && str.ToValue != "" {
			where += " and " + pf + "dec" + strconv.Itoa(key) + " between '" + str.FromValue + "' and '" + str.ToValue + "' "
		}
	}

	where += GetExtractStr(y, sqlStr)
	return where
}

// The function `GetLt1InsertProcFunc` parses a SQL string to extract and format a stored procedure call.
func GetLt1InsertProcFunc(y *abango.Controller, sql string) string {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Scan() //@procedure skip
	scanner.Scan()
	fnStr := "call " + scanner.Text()
	fmt.Println(e.FuncNameInfo() + ":" + fnStr)
	return fnStr
}

func GetLt1ProcFunc(y *abango.Controller, sql string, lt *ListType1Vars, ltFilter string) string {
	if ltFilter == "" { //Query 에러가 나지 않도록 하기 위해서 임.
		ltFilter = "0"
	} // ListTypeFilter 는 예비용 필터이며 and/or 또는 = 같은 SQL keyword가 들어가면 안된다.

	// 특수한 경우의 커스텀 파라메터-시작
	var storageIdStr string // 다중창고 관리시 특정 샃고 지정헤서 수불장 보기
	if !strings.Contains(lt.ListSimpleFilter, "@procedure.storage_id:") {
		storageIdStr = strconv.Itoa(y.Gtb.StorageId)
	} else {
		storageIdStr = strings.TrimPrefix(lt.ListSimpleFilter, "@procedure.storage_id:")
	}
	// 특수한 경우의 커스텀 파라메터-시작

	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Scan() //@procedure skip
	scanner.Scan()
	fnStr := "call " + scanner.Text()
	fnStr += "('" + lt.ListToken + "', '" + lt.StartDate + "', '" + lt.EndDate
	fnStr += "', '" + lt.StartFirst + "', '" + lt.EndFirst
	fnStr += "', '" + lt.StartSecond + "', '" + lt.EndSecond
	fnStr += "', '" + lt.StartThird + "', '" + lt.EndThird
	fnStr += "', '" + lt.StartFourth + "', '" + lt.EndFourth
	fnStr += "', '" + lt.Balance + "', '" + lt.OrderBy
	fnStr += "', " + strconv.Itoa(y.Gtb.BranchId) + ", " +
		storageIdStr + ", " + strconv.Itoa(y.Gtb.MemberCompanyId) + ", " + ltFilter
	// strconv.Itoa(y.Gtb.StorageId) + ", " + strconv.Itoa(y.Gtb.MemberCompanyId) + ", " + ltFilter
	fnStr += ")"
	fmt.Println(e.FuncNameInfo() + ":" + fnStr)
	return fnStr
}

func GetPopWhereCnd(p *PopupList1Vars) string {
	where := ""
	if p.SumFilterName != "" {
		where += " and " + p.SumFilterName + " = " + p.SumFilterValue + " "
	}

	if p.SumSimpleFilter != "" {
		where += " and " + p.SumSimpleFilter + " "
	}
	return where
}

func GetPvWhereCnd(pv *PageVars) string {
	var where string
	if pv.Query != "" {
		where += " and " + pv.Query + " "
	}
	return where
}

func GetQryWhereCnd(q *QueryVars) string {
	var where string
	if q.FilterName != "" && q.FilterValue != "" {
		where += " and " + q.FilterName + " like '%" + q.FilterValue + "%' "
		// int가 들어와도 형편환하여 like 검색 가능, right like 속도가 안떨어진다.
	}
	if q.SimpleFilter != "" {
		where += " and " + q.SimpleFilter + " "
	}
	return where
}

func Lt1QvWhereGet(y *abango.Controller, q *QueryVars, lt *ListType1Vars, sqlStr string) string {

	qwhereClause := GetQryWhereCnd(q)
	// FilterThird 는 항상 품목으로 하여 subfilter 로 범위지정에 되어있는 경우만 수행함.
	if lt.FilterThird == "item_code" && lt.StartThird != "" {
		q.SubSimpleFilter = " " + lt.FilterThird + " between '" + lt.StartThird + "' and '" + lt.EndThird + "' "
		lt.FilterThird = "" //FilterThird 는 뺀다.
	}
	return qwhereClause + GetLt1WhereCnd(y, lt, sqlStr)
}

func GetLt1WhereCnd(y *abango.Controller, lt *ListType1Vars, sqlStr string) string {
	var where string

	// ListFilterName + Value 처리
	if lt.ListFilterName != "" {
		column := lt.ListFilterName
		value := lt.ListFilterValue

		if strings.HasPrefix(column, "N_") {
			column = strings.TrimPrefix(column, "N_")
			where += fmt.Sprintf(" AND %s = '%s' ", column, value)
		} else {
			where += fmt.Sprintf(" AND %s LIKE '%%%s%%' ", column, value)
		}
	}

	// 간단 필터
	if lt.ListSimpleFilter != "" {
		where += fmt.Sprintf(" AND %s ", lt.ListSimpleFilter)
	}

	// 날짜/범위 필터 묶음
	rangeFilters := []struct {
		column string
		start  string
		end    string
	}{
		{lt.FilterDate, lt.StartDate, lt.EndDate},
		{lt.FilterFirst, lt.StartFirst, lt.EndFirst},
		{lt.FilterSecond, lt.StartSecond, lt.EndSecond},
		{lt.FilterThird, lt.StartThird, lt.EndThird},
		{lt.FilterFourth, lt.StartFourth, lt.EndFourth},
	}

	for _, rf := range rangeFilters {
		if rf.column != "" && rf.start != "" && rf.end != "" {
			where += fmt.Sprintf(" AND %s BETWEEN '%s' AND '%s' ", rf.column, rf.start, rf.end)
		}
	}

	where += GetExtractStr(y, sqlStr)
	return where
}

func GetBodyCopyWhereCnd(y *abango.Controller, bc *BodyCopyPageVars, sqlStr string) string {
	var whereBuilder strings.Builder

	addCondition := func(condition string) {
		if condition != "" {
			whereBuilder.WriteString(" and ")
			whereBuilder.WriteString(condition)
		}
	}

	// 조건별 where절 추가
	if bc.SlipNo != "" {
		addCondition(fmt.Sprintf("%s like '%%%s%%'", bc.SlipNoField, bc.SlipNo))
	}

	if bc.CompanyName != "" {
		addCondition(fmt.Sprintf("company_name like '%%%s%%'", bc.CompanyName))
	}

	if bc.ItemCode != "" {
		addCondition(fmt.Sprintf("item_code like '%%%s%%'", bc.ItemCode))
	}

	if bc.ShowOnlyClosed != "" {
		field := GetQueryCommentStr(sqlStr, QcClosed, QcEnd)
		addCondition(fmt.Sprintf("%s = '%s'", field, bc.ShowOnlyClosed))
	}

	if bc.DaysFromToday != "" {
		addCondition(GetDeliveryStr(sqlStr, QcDelivery, QcEnd, 1000))
	}

	// 추가 where절 추출
	whereBuilder.WriteString(GetExtractStr(y, sqlStr))

	return whereBuilder.String()
}

func GetMyWhereCnd(y *abango.Controller, my string) string {
	var value string

	switch {
	case strings.Contains(my, "member"):
		value = e.NumToStr(y.Gtb.MemberId)
	case strings.Contains(my, "buyer"), strings.Contains(my, "company"):
		value = e.NumToStr(y.Gtb.MemberCompanyId)
	case strings.Contains(my, "user"):
		value = e.NumToStr(y.Gtb.UserId)
	default:
		return "" // 조건에 해당하지 않으면 빈 문자열 반환
	}

	return fmt.Sprintf(" and %s=%s ", my, value)
}

func GetAllQueryCnd(y *abango.Controller, sql string, cnt string, where string,
	having string, order string, subWhere string, limitoffset string) (string, string) {

	replaceClause := func(src, placeholder, defaultClause, clause string, isCount bool) (string, string) {
		fullClause := ""
		if clause == "" {
			fullClause = defaultClause
		} else {
			if placeholder == QcSubWhere || placeholder == QcHaving {
				fullClause = defaultClause + " and " + clause + "\n"
			} else {
				fullClause = defaultClause + clause
			}
		}
		src = strings.Replace(src, placeholder, fullClause, -1)
		if isCount {
			cnt = strings.Replace(cnt, placeholder, fullClause, -1)
		}
		return src, cnt
	}

	sql, cnt = replaceClause(sql, QcWhere, QkWhere, where, cnt != "")
	sql, cnt = replaceClause(sql, QcSubWhere, QkWhere, subWhere, cnt != "")
	sql, cnt = replaceClause(sql, QcHaving, QkHaving, having, cnt != "")

	if order != "" {
		sql = strings.Replace(sql, QcOrder, QkOrder+order, -1)
	}
	if limitoffset != "" {
		sql = strings.Replace(sql, QcLimitOffset, limitoffset, -1)
	}

	return sql, cnt
}

func GetExtractStr(y *abango.Controller, sqlStr string) string {
	columns := strings.Split(GetQueryCommentStr(sqlStr, QcExtract, QcEnd), ",")
	var builder strings.Builder

	for _, col := range columns {
		col = strings.TrimSpace(col) // 혹시 모를 공백 제거

		switch {
		case strings.Contains(col, "user_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.UserId))

		case strings.Contains(col, "branch_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.BranchId))

		case strings.Contains(col, "member_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.MemberId))

		case strings.Contains(col, "company_id"),
			strings.Contains(col, "buyer_id"),
			strings.Contains(col, "seller_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.MemberCompanyId))

		case strings.Contains(col, "storage_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.StorageId))

		case strings.Contains(col, "store_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.StoreId))

		case strings.Contains(col, "terminal_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.TerminalId))

		case strings.Contains(col, "store_buyer_id"):
			builder.WriteString(fmt.Sprintf(" and %s=%d", col, y.Gtb.StoreBuyerId))

		}
	}

	return builder.String()
}

func GetQueryCommentStr(sqlStr string, start string, end string) string {

	posFirst := strings.Index(sqlStr, start)
	if posFirst == -1 { // 찾는 스트링이 없으면 그냥 빠져나간다.
		return ""
	}
	posFirst = posFirst + len(start)
	posLast := strings.Index(sqlStr[posFirst:], end) + posFirst
	if posLast == -1 {
		e.FuncRunErr("023jorfjs3w", e.FuncNameErr()+"'"+QcEnd+"' string NOT found in query string ! ")
		return ""
	}
	return sqlStr[posFirst:posLast] + " "
}

func GetDeliveryStr(sqlStr string, start string, end string, day int) string {
	deliveryDate := time.Now().UTC().AddDate(0, 0, day).Format("20060102")
	return GetQueryCommentStr(sqlStr, start, end) + "<= '" + deliveryDate + "' "
}

// 이건 abaongo로 갈 것 !!!
func FileListToSlice(path string) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, e.LogErr("3fdgfvaqrf", e.FuncNameErr(), err)
	}

	// Convert DirEntry slice to FileInfo slice
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, e.LogErr("3fdgfvaq", e.FuncNameErr(), err)
		}
		files = append(files, info)
	}

	return files, nil
}
