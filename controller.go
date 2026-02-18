package abango

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	e "github.com/dabory/abango-rest/etc"
	"github.com/go-xorm/xorm"
	"gopkg.in/yaml.v2"
)

func (c *Controller) Init(r *http.Request) (int, string) {

	// fmt.Println("Heare is abango.Controller.Init") // Just make it Comment line. DONOT Delete.

	var gtb GateTokenBase

	var gtbStr string
	var err error
	if XConfig["IsYDBFixed"] == "Yes" {

		data, err := os.ReadFile("models/custom.yml")
		if err != nil {
			return 507, e.LogStr("ASDFQEWFA", "Can NOT Read custom.yml")
		}

		// fmt.Println("aaa:"), string(data))

		var config Config
		if err := yaml.Unmarshal(data, &config); err == nil {
			c.Gtb = config.Source // 여기서 custom.yml 값이 할당이 안되면 'yaml: "????"  갑이 제대로 할당이 안된거다.
		} else {
			return 507, e.LogStr("ASDWEWFA", "connString in custom.yml format mismatch ")
		}

	} else {
		if c.GateToken == "" {
			return 505, e.LogStr("QWCAFVD", "GateToken is Empty: ")
		}
		if gtbStr, err = AegisView(r, c.GateToken); err != nil { // 여기의 DB접속시 r의 웹로그정보를 AegisCache에 전달.
			return 505, e.LogStr("QWFAECD", "GateToken Not Found in AegisDB: "+c.GateToken)
		}

		if err := json.Unmarshal([]byte(gtbStr), &gtb); err == nil {
			c.Gtb = gtb
		} else {
			return 505, e.LogStr("QWFAEC1AFVDS", "AfterBase64Content Format mismatch: "+c.GateToken)
		}
	}

	if status, msg := c.AttachDB(); status != 200 { // DB 까지 붙여야 memory error 가 안난다.
		return status, e.LogStr("PBUYJM-", msg)
	}
	return 200, ""
}

type Config struct {
	Source GateTokenBase
}

type Target struct {
	Type      string `yaml:"type"`
	Language  string `yaml:"language"`
	OutputDir string `yaml:"output_dir"`
}

func (c *Controller) AttachDB() (int, string) {

	var err error
	if c.Db, err = xorm.NewEngine(XConfig["DbType"], c.Gtb.ConnString); err != nil {
		return 600, e.LogStr("ADASEF", "DBEngine Open Error")
	}

	var connHint string
	strArr := strings.Split(c.Gtb.ConnString, "@tcp")
	if len(strArr) == 2 {
		connHint = strArr[1]
	} else {
		return 507, e.LogStr("ASDFQEWFA", "connString format mismatch: "+strArr[1])
	}

	c.Db.ShowSQL(false)
	c.Db.SetMaxOpenConns(100)
	c.Db.SetMaxIdleConns(20)
	c.Db.SetConnMaxLifetime(60 * time.Second)
	if _, err := c.Db.IsTableExist("aaa"); err == nil {
		return 200, e.LogStr("ASDFASFQFE", "YDB connection in "+connHint)
	} else {
		return 600, e.LogStr("PMUHIUYBUYJM-", "YDB connection Fail in "+connHint)
	}

	return 200, ""
}
