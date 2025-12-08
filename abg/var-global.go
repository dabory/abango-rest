package abg

import "github.com/dabory/abango-rest"

const ( // Query
	// Tpf string = "dbr_" // TablePrefix
	//QueryComment
	QcWhere        string = "-- @where"
	QcSubWhere     string = "-- @subwhere"
	QcHaving       string = "-- @having"
	QcOrder        string = "-- @order"
	QcLimitOffset  string = "-- @limitoffset"
	QcExtract      string = "-- @extract:"
	QcClosed       string = "-- @closed:"
	QcDelivery     string = "-- @delivery:"
	QcBetweenDates string = "-- @between_dates"
	QcEnd          string = "--" //QueryComment
	//QueryKeyword
	QkWhere string = "\nwhere true "
	// QkWhere  string = "\nwhere 1 "
	QkHaving string = "\nhaving true "
	QkOrder  string = "\norder by "
	QkLimit  string = "\nlimit  "
	QkOffset string = " offset "

	// 합계금액, order_by, t_id 이 순서로 해야만 한다.
	QkTmpOrder string = " order by is_sum desc, order_by, t_id asc "
	BAR        string = "@_@"
)

// const ( //SsoSubId

// 	KkCrm_SsoSubId int = 3
// )

type AppApi struct {
	BackUrl   string
	GateToken string
}

// 0:Sso, 1:Dbu  매우 중요하다.
var AppApis [2]AppApi

var ( //Env XCongif
	SQL_DEBUG          bool
	NORMAL_DEBUG       bool
	DEVICE_AUTH        bool
	DBU_BY_FORCE       bool
	IS_CACHE_KEY_PAIR  bool
	KAFKA_CONSUMER     bool
	LOCAL_KEY_PAIR     string
	CACHE_KEY_PAIR_DIR string
	QHOME_DIR          string
	THEME_QRY_DIR      string

	WEB_LOG_SW       string
	IS_LOG_DEBOUNCED string

	MAIN_PRODUCER_TOPIC string
	COMSUMER_TOPICS     []string
	CRY_ABG             abango.Controller
	Last                QryName
	// REDIS_EXTIME        time.Duration = 12 * time.Hour
)

type QryName struct {
	ListType1Page string
	ListType1Book string
}

var ( //queries dir
	COPY_DIR   string = "/copy"
	CHART_DIR  string = "/chart"
	FORM_DIR   string = "/form"
	FUNC_DIR   string = "/func"
	LIST_DIR   string = "/list"
	POPUP_DIR  string = "/popup"
	SEARCH_DIR string = "/search"
	TURBO_DIR  string = "/turbo"
)

var ( //Prefix
	ES_PREFIX string = "erp_"
)
