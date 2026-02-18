package abg

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/dabory/abango-rest"
	e "github.com/dabory/abango-rest/etc"
)

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
	LANG_DIR           string

	WEB_LOG_SW       string
	IS_LOG_DEBOUNCED string

	MAIN_PRODUCER_TOPIC string
	COMSUMER_TOPICS     []string
	CRY_ABG             abango.Controller
	Last                QryName
	// REDIS_EXTIME        time.Duration = 12 * time.Hour
)

var ( //InfluxDB
	INFLUX_CONN      string
	INFLUX_DB_NAME   string
	INFLUX_USER_NAME string
	INFLUX_PASSWORD  string
	INFLUX_DURATION  string
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

// abg -> abango, abg -> etc, abango->etc 이런 형태로 abg 가 모든 것을 호출할 수 있은 구조로 되어 있다.
// abango 는 주요 번수 설정을 위한 기반을 제공하고 그외에 abg 가 상세 설정을 하는 구조이다.수정하려고 하지 말것.
func GlobalVarsInit() error { // MainApi, AegisCache, StrongApi 모두 이걸 사용한다.

	abango.CpuCores = runtime.NumCPU()

	if err := abango.InitializeOTP(); err != nil {
		e.LogErr("OTP_INIT", "failed to initialize OTP manager", err)
		panic("")
	}

	abango.AEGIS_CACHE_CONN = abango.XConfig["AegisCacheConn"]
	if abango.AEGIS_CACHE_CONN == "" { // 할당하지 않을 경우 일반적으로 자신이 서버에서 동작된다고 간주
		abango.AEGIS_CACHE_CONN = "localhost:19090"
	} else { // AegisCacheConn 변수가 할당된 경우만 에러 테크한다. // 이건 특수한 경우이다.
		var connRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+:[0-9]+$`)
		if abango.AEGIS_CACHE_CONN == "" || !connRegex.MatchString(abango.AEGIS_CACHE_CONN) {
			e.LogErr("abango.AEGIS_CACHE_CONN: ", abango.AEGIS_CACHE_CONN, errors.New("is not a valid connection format (host:port)"))
			panic("Invalid connection format: " + abango.AEGIS_CACHE_CONN)
		}
	}

	// Global Variable Init - 극히 반복적으로 쓰이는 XConfig 변수는 Global로 만든다.
	DEVICE_AUTH = e.YesToTrue(abango.XConfig["DeviceAuthOn"], "DeviceAuthOn")
	DBU_BY_FORCE = e.YesToTrue(abango.XConfig["DbuByForceOn"], "DbuByForceOn")
	IS_CACHE_KEY_PAIR = e.YesToTrue(abango.XConfig["IsCacheKeyPair"], "IsCacheKeyPair")
	KAFKA_CONSUMER = e.YesToTrue(abango.XConfig["IsKafkaConsumer"], "IsKafkaConsumer")

	// 회원 민감정보 암호화 스위치
	abango.AEGIS_MEMBER_ON = e.YesToTrue(abango.XConfig["AegisMemberOn"], "AegisMemberOn")

	COMSUMER_TOPICS = strings.Split(strings.Replace(abango.XConfig["ConsumerTopics"], " ", "", -1), ",")
	MAIN_PRODUCER_TOPIC = abango.XConfig["MainProducerTopic"]
	LOCAL_KEY_PAIR = abango.XConfig["LocalKeyPair"]
	CACHE_KEY_PAIR_DIR = abango.XConfig["CacheKeyPairDir"]
	QHOME_DIR = abango.XConfig["QueryDir"]
	LANG_DIR = abango.XConfig["LangDir"]

	INFLUX_CONN = abango.XConfig["InfluxConn"]
	INFLUX_DB_NAME = abango.XConfig["InfluxDbName"]
	INFLUX_USER_NAME = abango.XConfig["InfluxUserName"]
	INFLUX_PASSWORD = abango.XConfig["InfluxPassword"]
	INFLUX_DURATION = abango.XConfig["InfluxDuration"]

	fmt.Println("===== InfluxDB Config =====")
	fmt.Println("INFLUX_CONN      :", INFLUX_CONN)
	fmt.Println("INFLUX_DB_NAME   :", INFLUX_DB_NAME)
	fmt.Println("INFLUX_USER_NAME :", INFLUX_USER_NAME)
	fmt.Println("INFLUX_PASSWORD :", INFLUX_PASSWORD) // 보안상 마스킹
	fmt.Println("INFLUX_DURATION :", INFLUX_DURATION) // 보안상 마스킹
	fmt.Println("===========================")

	// 이것은 mainapi 와 strong api 각각 수정해서 적용해야 함. theme 에서 단순 쿼리는 mainapi conf 에서 string api 용궈리는 strong conf 폴더에서 수정함.
	THEME_QRY_DIR = abango.XConfig["ThemeQryDir"]
	WEB_LOG_SW = abango.XConfig["WebLogSw"]
	IS_LOG_DEBOUNCED = abango.XConfig["IsLogDebounced"]

	return nil
}
