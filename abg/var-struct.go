// Author : Eric Kim
// Build Date : 23 Jul 2008  Last Update 02 Aug 2008
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package abg

// 0. Global Variable /////////

type ListFormBalanceVars struct {
	IsApplyBalance bool
	YyyyMm         string
	SelectedId     int
	StartCode      string
	EndCode        string
	OrderBy        string
}

type CompanySearchVars struct {
	CompanyName string
	MainContact string
	MobileNo    string
	Email       string
	OrderBy     string
}

type MediaSearchVars struct {
	StartDate    string
	EndDate      string
	SlipNo       string
	MediaName    string
	FileName     string
	LinkLocation string
	Linked       string
	NickName     string
	BranchName   string
	OrderBy      string
}

type ItemSearchVars struct {
	ItemCode string
	ItemName string
	SubName  string
	OrderBy  string
}

type SlipSearchVars struct {
	StartDate   string
	EndDate     string
	SlipNo      string
	CompanyName string
	ItemCode    string
	QuerySpeed  string
	OrderBy     string
}

type SlipSearchFields struct {
	SlipDateField string
	SlipNoField   string
}

type SingleVars struct {
	QueryName string
	Id        int
}

type Bb64Base struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Database string
	Password string
}

type SetupBase struct {
	SetupCode string
	SetupJson string
}

type QueryVars struct {
	QueryName       string
	FilterName      string
	FilterValue     string
	SimpleFilter    string
	SubSimpleFilter string
	IsntPagination  bool
	TestMode        string
}
type CopyTableVars struct {
	SourceTblName string
	TaegetTblName string
	SourceNoName  string
	TargetNoName  string
}

type SelectType1Vars struct {
	ListToken string
	NameValue string
	// Str            []FilterBase
	// Chk            []OptBase
	// Rng            []RngBase
	// Dec            []RngBase
	// StrItem        []FilterBase
	// ChkItem        []OptBase
	// RngItem        []RngBase
	// DecItem        []RngBase
	Having         SelectFilters
	Where          SelectFilters
	IsDownloadList bool
	Balance        string
	OrderBy        string
}

type SelectFilters struct {
	Prefix string
	Str    []FilterBase
	Chk    []OptBase
	Rng    []RngBase
	Dec    []RngBase
}

type OptBase struct {
	Opt []FilterBase
}

type FilterBase struct {
	FilterValue string
}

type RngBase struct {
	FromValue string
	ToValue   string
}

type PopupList1Vars struct {
	PopupFilterName  string
	PopupFilterValue string
	SumFilterType    string // "" : Body/Sum 동시적용 1: Sum만 적용, 7: Filter OR SumFilter 로 둘다 적용, 8: OR 로 묶어 Sum만 적용
	SumFilterName    string
	SumFilterValue   string
	SumSimpleFilter  string
}

type ListType1Vars struct {
	ListToken string

	FilterDate string
	StartDate  string
	EndDate    string

	FilterFirst string
	StartFirst  string
	EndFirst    string

	FilterSecond string
	StartSecond  string
	EndSecond    string

	FilterThird string
	StartThird  string
	EndThird    string

	FilterFourth string
	StartFourth  string
	EndFourth    string

	IsAddTotalLine   bool
	IsExcelColumn    bool
	IsShowOnlyClosed bool
	IsDownloadList   bool
	IsntPageReturn   bool
	IsTmpSqlUsed     bool
	IsCrystalReport  bool
	Balance          string
	OrderBy          string

	ListFilterName   string
	ListFilterValue  string
	ListSimpleFilter string
}

type BodyCopyPageVars struct {
	SlipNoField    string
	SlipNo         string
	ItemCode       string
	CompanyName    string
	ShowOnlyClosed string
	Balance        string
	DaysFromToday  string
	OrderBy        string
}

type PageVars struct {
	MyFilter   string
	QueryCnt   int64
	Query      string
	Fields     string
	Asc        string
	Desc       string
	Limit      int
	Offset     int
	ReturnJson string
}

// type InputFieldJson struct {
// 	StartDate  string
// 	EndDate    string
// 	Order      int
// 	Select     int
// 	QueryInput string
// }

// type Login struct {
// 	UserId   string
// 	Password string
// }

// // 3. Answer /////////////////////////////////////////////////////////////////
// type AnswerBase struct {
// 	SvcSts string
// 	SvcMsg string
// }

// type ActRowCom struct {
// 	AnswerBase
// 	IdBase
// }

// type ActPageCom struct {
// 	AnswerBase
// 	IdPageBase
// }

type BodyCopyBase struct {
	BdId int
	Qty  string
}

type IdBase struct {
	Id int
}

type SimpleQryBase struct {
	Idx       int // MainId
	Idh       int // HeadId
	Ids       int // SubId
	CreatedOn int64
	UpdatedOn int64
	C1        string
	C2        string
	C3        string
	C4        string
	C5        string
	C6        string
	C7        string
	C8        string
	C9        string
	C10       string
	OrderBy   string
}

type IdPageBase struct {
	Page []IdBase
}

// type AbangoApp struct {
// 	YDB *xorm.Engine
// }

type InsertVars struct {
	QueryName       string
	InsertType      string
	ListToken       string
	PreProcess      string
	PostProcess     string
	IsTruncateTable bool
	IsBackupTable   bool
	IsBackupDb      bool
}

type UploadBatchVars struct {
	UploadBatch   string
	IsCreateMedia bool
	IsCropImage   bool
	ImgExtension  string
}

// type ElasticInsertVars struct {
// 	Url      string
// 	Index    string
// 	Username string
// 	Password string
// }

//	type GateTokenGetReq struct {
//		// Target        string
//		ClientId     string
//		BeforeBase64 string
//		OwnerKey     string
//		// 아래는 deprecate 예정
//		AppBase64     string
//		Api23Key      string
//		Api23eKeyPair string
//	}
//
//	type KeyPairGetReq struct {
//		ClientId string
//	}
type DummyReq struct {
}

type IsMymenuSetReq struct {
	TableCode string
	MenuId    int
	IsMymenu  string
}
type MemberAuthCom struct {
	// member
	Email        string
	Password     string
	FirstName    string
	SurName      string
	NickName     string
	ActivateCode string
	SsoBrand     string
	SsoSub       string

	// member_ext
	MobileNo        string `xorm:"VARCHAR(20)"`
	PhoneNo         string `xorm:"VARCHAR(20)"`
	Sex             string `xorm:"not null default 'm' comment('m:남성, w:여성') CHAR(1)"`
	BirthDate       string `xorm:"not null default '' comment('생년월일') CHAR(8)"`
	Ssn1            string `xorm:"not null default '' comment('주민번호 앞자리') VARCHAR(10)"`
	Ssn2            string `xorm:"not null default '' comment('주민번호 뒷자리') VARCHAR(10)"`
	SellerCode      string `xorm:"not null default '' comment('판매자코드') unique VARCHAR(32)"`
	IsSeller        string `xorm:"not null default '0' comment('판매자승인완료(판매자임)') CHAR(1)"`
	IsMyappOk       string `xorm:"not null default '0' comment('마이앱 사용가능') CHAR(1)"`
	SellerConfirmOn int64  `xorm:"not null default 0 comment('판매자 신청일시') BIGINT(20)"`
	SellerRequestOn int64  `xorm:"not null default 0 comment('판매자 신청일시') BIGINT(20)"`
	SsnCardImg      string `xorm:"not null default '' comment('주민증 이미지') VARCHAR(256)"`
	MemberPermId    int    `xorm:"not null INT(11)"`
	MenuLangSw      int    `xorm:"not null TINYINT(4)"`
	SgroupId        int    `xorm:"not null INT(11)"`
	BranchId        int    `xorm:"not null INT(11)"`
	StorageId       int    `xorm:"not null INT(11)"`
	AgroupId        int    `xorm:"not null INT(11)"`
	CountryCode     string `xorm:"not null CHAR(5)"`
	IsExpired       string `xorm:"not null default '0' CHAR(1)"`
	CustomA         string `xorm:"not null default '' comment('커스텀용 필드 A') VARCHAR(64)"`
	CustomB         string `xorm:"not null default '' comment('커스텀용 필드 A') VARCHAR(64)"`
	CustomC         string `xorm:"not null default '' comment('커스텀용 필드 A') VARCHAR(64)"`
	CustomD         string `xorm:"not null default '' comment('커스텀용 필드 A') VARCHAR(64)"`
	CustomE         string `xorm:"not null default '' comment('커스텀용 필드 A') VARCHAR(64)"`

	// DbrCompany 중복되는 필드는 Email, MobileNo, Sex, BirthDate
	CompanyDate       string `xorm:"not null default '20200101' comment('업체등록일') index CHAR(8)"`
	CompanyNo         string `xorm:"not null default '20220101-01' comment('업체등록 번호') unique VARCHAR(21)"`
	CgroupId          int    `xorm:"comment('업체구분') index INT(10)"`
	SellerId          int    `xorm:"not null default 1 comment('상위판매자:Default 1') INT(11)"`
	Sort              string `xorm:"not null default 'AA' comment('고객sort') CHAR(2)"`
	CompanyName       string `xorm:"comment('고객약칭') index VARCHAR(32)"`
	CompanyClass      string `xorm:"default '' comment('업체등급:AA(개인고객), AB(기업고객), BB(기업공급처)') index VARCHAR(2)"`
	FullName          string `xorm:"not null default '' comment('고객성명(전체이름)') VARCHAR(96)"`
	IsLunar           string `xorm:"comment('0:양력, 1:음력') CHAR(1)"`
	CardChar4         string `xorm:"comment('지울것') VARCHAR(4)"`
	MainContact       string `xorm:"comment('주담당자: POS의 경우 고객이름(동명이인 사용가능)으로 사용') VARCHAR(21)"`
	TelNo             string `xorm:"index VARCHAR(21)"`
	FaxNo             string `xorm:"VARCHAR(21)"`
	TaxNo             string `xorm:"comment('사업자 등록 번호') VARCHAR(21)"`
	President         string `xorm:"comment('대표자명') VARCHAR(96)"`
	ZipCode           string `xorm:"comment('우편번호') VARCHAR(21)"`
	Addr1             string `xorm:"comment('현주소') VARCHAR(49)"`
	Addr2             string `xorm:"comment('상세주소') VARCHAR(49)"`
	BizType           string `xorm:"comment('업태') VARCHAR(191)"`
	DealItem          string `xorm:"comment('종목') VARCHAR(191)"`
	IsDealEnd         string `xorm:"not null default '0' comment('거래중지') CHAR(1)"`
	IsOkText          string `xorm:"not null default '0' comment('문자접수승인') CHAR(1)"`
	IsOkEmail         string `xorm:"not null default '0' comment('이메일접수승인') CHAR(1)"`
	IsOkDm            string `xorm:"default '0' comment('DM접수승인') CHAR(1)"`
	CurrCreditBal     string `xorm:"not null default 0.0000 comment('지울것') DECIMAL(20,4)"`
	CourierCode       string `xorm:"not null default '' comment('지정 택배사코드(etc에 있슴)') VARCHAR(8)"`
	Remarks           string `xorm:"MEDIUMTEXT"`
	CertImg           string `xorm:"default '' comment('사업자등록증 사진') VARCHAR(256)"`
	OnlineCertNo      string `xorm:"not null default '' comment('세금계산서이메일') VARCHAR(64)"`
	TaxMail           string `xorm:"not null default '' comment('세금계산서이메일') VARCHAR(64)"`
	ShopAbroad        string `xorm:"not null default '0' comment('사업자위치:0국내,1해외') CHAR(1)"`
	ShopName          string `xorm:"not null default '온라인상점명' comment('온라인상점명') VARCHAR(64)"`
	OnlineCertImg     string `xorm:"not null default '' comment('통신판매신고증') VARCHAR(256)"`
	ShipType          string `xorm:"not null default '0' comment('배송타입;0오늘출발') VARCHAR(21)"`
	ShipFeeBrand      string `xorm:"not null default 'free' comment('배송비setup-brand') VARCHAR(32)"`
	ReturnFee         string `xorm:"not null default 0.0000 comment('반품배송비') DECIMAL(18,4)"`
	AvgDeliDays       int    `xorm:"not null default 3 comment('평균배송기간(일)') INT(11)"`
	ExchangeFee       string `xorm:"not null default 0.0000 comment('교환배송비') DECIMAL(18,4)"`
	ShipZip           string `xorm:"not null default '999-999' comment('발송지우편번호') VARCHAR(12)"`
	ShipAddr1         string `xorm:"not null default '발송지주소1' comment('발송지주소1') VARCHAR(49)"`
	ShipAddr2         string `xorm:"not null default '발송지주소2' comment('발송지주소2') VARCHAR(49)"`
	ReturnZip         string `xorm:"not null default '999-999' comment('반품지우편번호') VARCHAR(12)"`
	ReturnAddr1       string `xorm:"not null default '반품지주소1' comment('반품지주소1') VARCHAR(49)"`
	ReturnAddr2       string `xorm:"not null default '반품지주소2' comment('반품지주소2') VARCHAR(49)"`
	SellerSettleBrand string `xorm:"not null default 'domestic-a' comment('판매자정산setup-brand') VARCHAR(32)"`
	CommissionRate    string `xorm:"not null default 0.0000 comment('판매수수료율') DECIMAL(10,4)"`
	NationCode        string `xorm:"not null default '' comment('국가코드') VARCHAR(5)"`
	CurrencyCode      string `xorm:"not null default '' comment('통화단위') VARCHAR(5)"`
	BankName          string `xorm:"not null default '' comment('은행명') VARCHAR(32)"`
	AccountNo         string `xorm:"not null default '' comment('입급계좌') VARCHAR(32)"`
	HolderName        string `xorm:"not null default '' comment('예금주') VARCHAR(32)"`
	AccountImg        string `xorm:"not null default '' comment('통장사본 이미지') VARCHAR(256)"`
	SalesBrand        string `xorm:"not null default '' comment('영업브랜드명') VARCHAR(64)"`
	SiteUrl           string `xorm:"not null default '' comment('웹사이트 URL') VARCHAR(64)"`
	SnsAccount        string `xorm:"not null default '' comment('SNS 계정') VARCHAR(64)"`
	CompanyJson       string `xorm:"comment('커스터마이징 내용') TEXT"`
	Ip                string `xorm:"VARCHAR(21)"`

	// ProMemberDevice
	// DeviceIp   string
	// DeviceDesc string
}
type UserAuthCom struct {
	ActivateCode string
	SsoBrand     string
	SsoSub       string

	Email             string
	Password          string
	FirstName         string
	SurName           string
	NickName          string
	MobileNo          string
	IsSkipDbupdateEnv string
	DbupdateRange     string //Login 시 Update Range 를 줌
}

type SetupRowReq struct {
	SetupCode string
	BrandCode string
	LangType  string
}

type TextVars struct {
	Email         string
	Encrypted     string
	BrandCode     string
	TemplateCode  string
	Sender        string
	ReservedTime  string
	TemplateTitle string
	TemplateText  string
	UniqueImage   string
}

type TextPageBase struct {
	ReplaceVars []ReplaceBase
}

type ReplaceBase struct {
	VarName  string
	VarValue string
}
