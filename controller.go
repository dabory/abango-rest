package abango

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	e "github.com/dabory/abango-rest/etc"
	"github.com/go-xorm/xorm"
	"gopkg.in/yaml.v2"
)

func (c *Controller) Init() (int, string) {

	var gtb GateTokenBase

	var gtbStr string
	var err error
	if XConfig["IsYDBFixed"] == "Yes" {

		data, err := os.ReadFile("models/custom.yml")
		if err != nil {
			return 507, e.LogStr("ASDFQEWFA", "Can NOT Read custom.yml")
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err == nil {
			c.Gtb = config.Source
			c.Gtb.ConnString = config.Source.ConnStr // custom.yml의 Variable name 이 서로 달라서 복사해줌.
		} else {
			return 507, e.LogStr("ASDWEWFA", "connString in custom.yml format mismatch ")
		}

	} else {
		if c.GateToken == "" {
			return 505, e.LogStr("QWCAFVD", "GateToken is Empty: ")
		}
		if gtbStr, err = MdbView(c.GateToken); err != nil {
			return 505, e.LogStr("QWFAECD", "GateToken Not Found in MemoryDB: "+c.GateToken)
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

func (c *Controller) CustomAbangoGet(ymlPath string) (int, string) {

	var err error
	data, err := os.ReadFile(ymlPath)
	if err != nil {
		return 507, e.LogStr("ASDFQEWFA", "Can NOT Read custom.yml")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err == nil {
		c.Gtb = config.Source
		c.Gtb.ConnString = config.Source.ConnStr // custom.yml의 Variable name 이 서로 달라서 복사해줌.
	} else {
		return 507, e.LogStr("ASDWEWFA", "connString in custom.yml format mismatch ")
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

// func (c *Controller) Init(ask AbangoAsk) {
// c.ServerVars = make(map[string]string) // 반드시 할당해줘야 함.
// c.Data = make(map[interface{}]interface{})

// c.Ctx.Ask = ask
// c.ConnString = XConfig["KafkaConnect"]
// c.Ctx.ReturnTopic = c.Ctx.Ask.UniqueId

// for _, p := range c.Ctx.Ask.ServerParams {
// 	c.ServerVars[p.Key] = p.Value
// }

// c.GetAbangoAccessAndDb()
// }

// func (c *Controller) InitNormal() {
// 	c.ServerVars = make(map[string]string) // 반드시 할당해줘야 함.
// 	c.Data = make(map[interface{}]interface{})
// 	//c.Ctx.ReturnTopic = "ljsldjfalsdfja" // 여기서 메모리 할당이 한됨'

// 	c.GetAbangoAccessAndDb()
// }

// func (c *Controller) KafkaAnswer(body string) {

// 	// c.Ctx.Answer.Body = []byte(body) // 쓸데없는 것 같은데 나중에 지
// 	// fmt.Println("ReturnTopic=" + c.Ctx.ReturnTopic)
// 	if _, _, err := KafkaProducer(body,
// 		c.Ctx.ReturnTopic, c.ConnString, XConfig["api_method"]); err != nil {
// 		e.MyErr("WERRWEEWQRFDFHQW", err, false)
// 	}
// }

func (c *Controller) AnswerJson() {

	// var ret []byte
	// if c.Data["json"] == nil {
	// 	msg := " QVCZEFAQWERQ " + c.Ctx.Ask.AskName + "[\"" + c.Ctx.Ask.ApiType + "\"] Data[json] is empty !"
	// 	e.MyErr(msg, errors.New(""), true)
	// 	return
	// }
	// if c.ServerVars["indent_answer"] == "yes" {
	// 	ret, _ = json.MarshalIndent(c.Data["json"], "", "  ")
	// } else {
	// 	ret, _ = json.Marshal(c.Data["json"])
	// }
	// // fmt.Println(string(ret))
	// if c.Ctx.Ask.ApiType == "Kafka" {
	// 	c.KafkaAnswer(string(ret))
	// } else if c.Ctx.Ask.ApiType == "gRpc" {
	// 	// 점진적으로 채워나가자.
	// } else if c.Ctx.Ask.ApiType == "Rest" {
	// 	// 점진적으로 채워나가자.
	// }

}

func (c *Controller) GetAbangoAccessAndDb() error {

	// if err := c.GetAccessAuth(); err == nil {
	// 	var err2 error
	// 	if c.Db, err2 = xorm.NewEngine(c.Access.DbType, c.Access.DbConnStr); err2 == nil {
	// 		// db, err := xorm.NewEngine(XEnc.DbType, "root:root@tcp(127.0.0.1:3306)/kangan?charset=utf8&parseTime=True")

	// 		strArr := strings.Split(c.Access.DbConnStr, "@tcp")
	// 		if len(strArr) == 2 {
	// 			e.LogNil(strArr[1])
	// 		} else {
	// 			e.MyErr(strArr[1], err2, true)
	// 			return err2
	// 		}

	// 		c.Db.ShowSQL(false)
	// 		c.Db.SetMaxOpenConns(100)
	// 		c.Db.SetMaxIdleConns(20)
	// 		c.Db.SetConnMaxLifetime(60 * time.Second)
	// 		if _, err := c.Db.IsTableExist("admin_menu"); err != nil {
	// 			e.MyErr("DB DISconnected in "+strArr[1], err, true)
	// 		} else {
	// 			e.LogNil("DB connect in " + strArr[1])
	// 		}
	// 	} else {
	// 		e.MyErr("xorm.NewEngine", err, true)
	// 	}
	// } else {
	// 	e.MyErr("GetAccessAuth", err, true)
	// }

	return nil
}

func (c *Controller) GetAccessAuth() error {

	return nil
}
