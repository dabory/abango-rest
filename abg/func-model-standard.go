package abg

import (
	"reflect"
	"strings"
	"time"

	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
	"github.com/google/uuid"
)

// 이거 다시 한번 원전히 정리 할 것.
func ComAddaRowDefault(y *abango.Controller, table interface{}) {
	fe := reflect.ValueOf(table).Elem()
	fieldNum := fe.NumField()
	for i := 0; i < fieldNum; i++ {
		field := fe.Type().Field(i)
		name := field.Name
		value := fe.Field(i)
		tag := string(field.Tag)

		switch {
		case strings.Contains(tag, "DECIMAL") && value.String() == "":
			fe.Field(i).SetString("0.0000")
		case name == "CreatedOn" || name == "RecordedOn":
			fe.Field(i).SetInt(e.GetNowUnix())
		case name == "OfficialDate":
			fe.Field(i).SetString(time.Now().Format("20060102"))
		case name == "OfficialTime":
			fe.Field(i).SetString(time.Now().Format("15:04:05"))
		case name == "IsUnused" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Status" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Sort" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Ip":
			fe.Field(i).SetString(y.Gtb.RemoteIp)
		case name == "UserId" && y.Gtb.UserId != 0:
			fe.Field(i).SetInt(int64(y.Gtb.UserId))
		case name == "MemberId" && y.Gtb.MemberId != 0:
			fe.Field(i).SetInt(int64(y.Gtb.MemberId))
		case name == "StorageId" && value.Int() == 0:
			fe.Field(i).SetInt(int64(y.Gtb.StorageId))
		case name == "BranchId" && value.Int() == 0:
			fe.Field(i).SetInt(int64(y.Gtb.BranchId))
		case name == "SgroupId" && value.Int() == 0:
			fe.Field(i).SetInt(int64(y.Gtb.SgroupId))
		case name == "AgroupId" && value.Int() == 0:
			fe.Field(i).SetInt(int64(y.Gtb.AgroupId))
		case name == "MemberBuyerId" && y.Gtb.MemberCompanyId != 0:
			fe.Field(i).SetInt(int64(y.Gtb.MemberCompanyId))
		case name == "BuyerId" && value.String() == "" && y.Gtb.MemberCompanyId != 0:
			fe.Field(i).SetInt(int64(y.Gtb.MemberCompanyId))
		case name == "MemberCompanyId" && y.Gtb.MemberCompanyId != 0:
			fe.Field(i).SetInt(int64(y.Gtb.MemberCompanyId))
		case name == "ItemCode" && value.String() == "":
			fe.Field(i).SetString(e.RandString(16))
		case name == "Duid" && value.String() == "":
			fe.Field(i).SetString(uuid.New().String())
		case (strings.Contains(name, "MediaId") || name == "MemberId" || name == "UserId" || name == "BuyerId" || name == "SupplierId" || name == "IgroupId" || name == "CgroupId" || name == "BranchId" || name == "StorageId" || name == "AgroupId") && value.Int() == 0:
			fe.Field(i).SetInt(int64(1))
		case name == "FromBuyerId" && value.Int() == 0:
			fe.Field(i).SetInt(int64(1))
		}
	}
}

func ComEditaRowDefault(y *abango.Controller, table interface{}) {
	fe := reflect.ValueOf(table).Elem()
	fieldNum := fe.NumField()

	for i := 0; i < fieldNum; i++ {
		field := fe.Type().Field(i)
		name := field.Name
		value := fe.Field(i)
		tag := string(field.Tag)

		switch {
		case strings.Contains(tag, "DECIMAL") && value.String() == "":
			fe.Field(i).SetString("0.0000")
		case name == "UpdatedOn" || name == "ModifiedOn":
			fe.Field(i).SetInt(e.GetNowUnix())
		case name == "OfficialDate" && value.String() == "":
			fe.Field(i).SetString(time.Now().Format("20060102"))
		case name == "OfficialTime" && value.String() == "":
			fe.Field(i).SetString(time.Now().Format("15:04:05"))
		case name == "UserId" && value.String() == "":
			fe.Field(i).SetInt(int64(y.Gtb.UserId))
		case name == "MemberId" && value.String() == "":
			fe.Field(i).SetInt(int64(y.Gtb.MemberId))
		case (name == "MemberBuyerId" || name == "MemberCompanyId") && value.String() == "":
			fe.Field(i).SetInt(int64(y.Gtb.MemberCompanyId))
		case name == "IsUnused" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Status" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Sort" && value.String() == "":
			fe.Field(i).SetString("0")
		case name == "Ip":
			fe.Field(i).SetString(y.Gtb.RemoteIp)
		case (strings.Contains(name, "MediaId") || name == "MemberId" || name == "UserId" || name == "BuyerId" || name == "SupplierId" || name == "IgroupId" || name == "CgroupId" || name == "BranchId" || name == "StorageId" || name == "AgroupId") && value.Int() == 0:
			fe.Field(i).SetInt(int64(1))
		}
	}
}
