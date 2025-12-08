package abg

// 여기는 모두 Array 형태로만 사용된다.

// Common Fixed Json
type GiftcardSettle struct { // Denomination Count - 권종별 갯수
	Sort  string // redeem:상품권회수, refund:환불
	Count int    // Sort별 횟수
	Amt   string // Sort별 금액
}
type CardSettle struct { // Denomination Count - 권종별 갯수
	Sort  string // approve:승인, cancel:취소
	Count int    // Sort별 횟수
	Amt   string // Sort별 금액
}
type DenomCount struct { // Denomination Count - 권종별 갯수
	Unit  string // 기타는 1 로 처리, 5000원 10000원 등
	Count int    // 권종별 갯수
}
type EmbededReply struct {
	CreatedOn  int64
	QnaSw      string // A:answer, B:questiongitpp
	WriterName string
	Contents   string
}

// DbrItem
type OptionAttr struct {
	Opt1   string
	Opt2   string
	Opt3   string
	Prc    string `xorm:"not null default 0.0000 DECIMAL(20,4)"`
	Qty    string `xorm:"not null default 0.0000 DECIMAL(20,4)"`
	Status bool
}

type ItemAddon struct {
	Caption string
	Direct  string
	Value   string
}

type ItemAttr struct {
	Caption string
	Value   string
}

type ItemBadge struct {
	BadgeName string
	BadgeType string
	CssClass  string
}

type ItemPrenext struct {
	Prenext    string
	Id         int
	TurboThumb string
	Title      string
	Slug       string
}

type ReviewIndex struct {
	Sort       string
	IndexCode  string
	Name       string
	LangType   string
	DeviceType string
	IsHidden   bool
}

// Post
type PostPrenext struct {
	Prenext    string
	Id         int
	TurboThumb string
	Title      string
	Slug       string
}

// ProItemReview
type ReviewScore struct {
	IndexCode string
	Score     int
}

// ProBalItemReview
type IndexSum struct {
	IndexCode string
	Score1Cnt int
	Score2Cnt int
	Score3Cnt int
	Score4Cnt int
	Score5Cnt int
}
