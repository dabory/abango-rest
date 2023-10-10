package spider

type WooTypeJson struct {
	Url
	Css
}

type Url struct {
	IsHtmlMatch  bool // false(defaul): Quick Match with Page Url, true: Page Html (ex:Magento)
	HtmlMatchStr string
	PdpRegEx     string //Product Detail Page Regular Expression
	PlpRegEx     string //Product List Page Regular Expression
}

type Css struct {
	//Css-Card1
	ItemName   string
	Categories string
	Currency   string
	SalesPrice string
	Images     string
	ShortDesc  string
	TextDesc   string
	Tags       string
	Emails     string
	//Css-Card2
	Language      string
	BrandName     string
	Manufacturer  string
	Sku           string
	Options       string
	UserCredit    string
	ModelNo       string
	ModelName     string
	Origin        string
	OriginDesc    string
	MinimumQty    string
	DeliveryPrice string
}
