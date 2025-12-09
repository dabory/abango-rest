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

// !!!여기를 수정할 경우 (1)UpdateGtb (2)MemberLoginSub (3)MemberSsoLogin (4)POSLogin 을 수정해야 한다.
type GateTokenBase struct {
	Password            string `yaml:"password"`
	ConnString          string `yaml:"conn_string"`
	AccessIp            string `yaml:"remote_ip"`
	UserAgent           string `yaml:"user_agent"`
	FrontIp             string `yaml:"front_ip"`
	FrontHost           string `yaml:"front_host"`
	Referer             string `yaml:"referer"`
	SsoSubId            int    `yaml:"sso_sub_id"`
	UserId              int    `yaml:"user_id"`
	UserPermId          int    `yaml:"user_perm_id"`
	UserFixedSortMenuId int    `yaml:"user_fixed_sort_menu_id"`
	UserSortMenuId      int    `yaml:"user_sort_menu_id"`

	MemberId  int    `yaml:"member_id"`
	NickName  string `yaml:"nick_name"`
	FirstName string `yaml:"first_name"`
	SurName   string `yaml:"sur_name"`

	MemberPermId          int    `yaml:"member_perm_id"`
	MemberFixedSortMenuId int    `yaml:"member_fixed_sort_menu_id"`
	SgroupId              int    `yaml:"sgroup_id"`
	BranchId              int    `yaml:"branch_id"`
	StorageId             int    `yaml:"storage_id"`
	AgroupId              int    `yaml:"agroup_id"`
	MemberCompanyId       int    `yaml:"member_company_id"`
	CompanySort           string `yaml:"company_sort"`
	SalesQtyPoint         int    `yaml:"sales_qty_point"`
	SalesPrcPoint         int    `yaml:"sales_prc_point"`
	SalesAmtPoint         int    `yaml:"sales_amt_point"`
	RetailQtyPoint        int    `yaml:"retail_qty_point"`
	RetailPrcPoint        int    `yaml:"retail_prc_point"`
	RetailAmtPoint        int    `yaml:"retail_amt_point"`
	PurchQtyPoint         int    `yaml:"purch_qty_point"`
	PurchPrcPoint         int    `yaml:"purch_prc_point"`
	PurchAmtPoint         int    `yaml:"purch_amt_point"`
	StockQtyPoint         int    `yaml:"stock_qty_point"`
	StockPrcPoint         int    `yaml:"stock_prc_point"`
	StockAmtPoint         int    `yaml:"stock_amt_point"`
	AccAmtPoint           int    `yaml:"acc_amt_point"`
	OfcCode               string `yaml:"ofc_code"`
	SalesVatSw            string `yaml:"sales_vat_sw"`
	RetailVatSw           string `yaml:"retail_vat_sw"`
	PurchVatSw            string `yaml:"purch_vat_sw"`
	StdVatRate            string `yaml:"std_vat_rate"`

	// POS 판매원 설정 //StorageId와 BranchId 는 이미 있슴
	TerminalId int `yaml:"terminal_id"`
	StoreId    int `yaml:"store_id"`

	//추가 할 경우 custom.yml 변수와 'yaml: "????" 과 일치가 필수입 !!!
}

type AgentTokenBase struct {
	ApiKey string
}
