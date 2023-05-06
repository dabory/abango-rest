package abango

import (
	"github.com/go-xorm/xorm"
)

var (
	XConfig   map[string]string
	FrontVars map[string]string //Fronrt End Server Variables
)

var (
	XDB *xorm.Engine
)

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

type GateTokenBase struct {
	ConnString    string
	RemoteIp      string
	DeviceDesc    string
	FrontIp       string
	FrontHost     string
	Referer       string
	SsoSubId      int
	UserId        int
	UserPermId    int
	MemberId      int
	MemberPermId  int
	SgroupId      int
	BranchId      int
	StorageId     int
	AgroupId      int
	MemberBuyerId int
	SalesQtyPoint int
	SalesPrcPoint int
	SalesAmtPoint int
	PurchQtyPoint int
	PurchPrcPoint int
	PurchAmtPoint int
	StockQtyPoint int
	StockPrcPoint int
	StockAmtPoint int
	AccAmtPoint   int
	OfcCode       string
}
