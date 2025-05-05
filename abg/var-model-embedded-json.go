package abg

// 여기는 모두 Array 형태로만 사용된다.

// Common Fixed Json
type EmbededReply struct {
	CreatedOn  int64
	QnaSw      string // A:answer, B:question
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
