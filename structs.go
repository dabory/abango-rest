package abango

import (
	"github.com/go-xorm/xorm"
)

var (
	XConfig   map[string]string
	FrontVars map[string]string //Fronrt End Server Variables
)

var (
	XDB       *xorm.Engine // 더 이상 쓰지 않음 2024에 없앨것
	CrystalDB *xorm.Engine // 더 이상 쓰지 않음 2024에 없앨것
)

// 1. Receivers /////////////////////////////////////////////////////////////////
// 1. Receivers /////////////////////////////////////////////////////////////////
type Param struct {
	Key   string
	Value string
}

type Controller struct {
	// Ctx            *context.Context
	Ctx            Context
	actionName     string
	controllerName string
	ServerVars     map[string]string //Fronrt End Server Variables
	GlobalVars     map[string]string //Fronrt End Global Variables
	Data           map[interface{}]interface{}
	Access         AbangoAccess

	GateToken       string
	DeviceHash      string
	UpdateFieldList string
	Gtb             GateTokenBase
	Db              *xorm.Engine
	V               interface{}
}

type Context struct {
	Ask         AbangoAsk
	Answer      AbangoAnswer
	ReturnTopic string
}

type AbangoAsk struct {
	ApiType      string
	AskName      string
	AccessToken  string
	UniqueId     string
	DocRoot      string
	Body         []byte
	ServerParams []Param
}

type AbangoAnswer struct {
	Body []byte
}

type AbangoAccess struct {
	UserId    int64
	UserGuid  string
	UserName  string
	NickName  string
	DbType    string
	DbConnStr string
}

// !!!여기를 수정할 경우 (1)UpdateGtb (2)MemberLogin (3)MemberSsoLogin 을 수정해야 한다.
type GateTokenBase struct {
	ConnStr         string `yaml:"conn_str"` //Only from Custom.yml
	ConnString      string //RealConnection String
	RemoteIp        string `yaml:"remote_ip"`
	DeviceDesc      string `yaml:"device_desc"`
	FrontIp         string `yaml:"front_ip"`
	FrontHost       string `yaml:"front_host"`
	Referer         string `yaml:"referer"`
	SsoSubId        int    `yaml:"sso_sub_id"`
	UserId          int    `yaml:"user_id"`
	UserPermId      int    `yaml:"user_perm_id"`
	MemberId        int    `yaml:"member_id"`
	MemberPermId    int    `yaml:"member_perm_id"`
	SgroupId        int    `yaml:"sgroup_id"`
	BranchId        int    `yaml:"branch_id"`
	StorageId       int    `yaml:"storage_id"`
	AgroupId        int    `yaml:"agroup_id"`
	MemberCompanyId int    `yaml:"member_company_id"`
	CompanySort     string `yaml:"company_sort"`
	SalesQtyPoint   int    `yaml:"sales_qty_point"`
	SalesPrcPoint   int    `yaml:"sales_prc_point"`
	SalesAmtPoint   int    `yaml:"sales_amt_point"`
	PurchQtyPoint   int    `yaml:"purch_qty_point"`
	PurchPrcPoint   int    `yaml:"purch_prc_point"`
	PurchAmtPoint   int    `yaml:"purch_amt_point"`
	StockQtyPoint   int    `yaml:"stock_qty_point"`
	StockPrcPoint   int    `yaml:"stock_prc_point"`
	StockAmtPoint   int    `yaml:"stock_amt_point"`
	AccAmtPoint     int    `yaml:"acc_amt_point"`
	OfcCode         string `yaml:"ofc_code"`
	SalesVatSw      string `yaml:"sales_vat_sw"`
	PurchVatSw      string `yaml:"purch_vat_sw"`
	StdVatRate      string `yaml:"std_vat_rate"`
}
