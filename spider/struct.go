package spider

type SchDitem struct {
	PdplinkHash string
	TopUrl      string
	TargetPath  string // URL without scheme info
	HtmlHash    string
	HubUrls     string
	IgroupCode  string
	Categories  string
	Emails      string
	ItemName    string
	Tags        string
	Currency    string
	SalesPrice  string

	Images    string
	ShortDesc string
	TextDesc  string

	OgMetas string
	Options []struct {
		Name    string
		Choices []struct {
			Name  string
			Price string
		}
	}

	OriginDesc    string
	Manufacturer  string
	Origin        string
	Language      string
	DeliveryPrice string
	Sku           string
	// ItemNick      string
	// ModelName     string
	// ModelNo       string
	// BrandName     string
	// MinimumQty    float32
	// UserCredit    float32
	// Suggest     []string
	// Cats        []string
}
