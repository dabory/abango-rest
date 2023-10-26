package spider

import (
	"encoding/json"
	"errors"

	e "github.com/dabory/abango-rest/etc"
	"github.com/elastic/go-elasticsearch/v7"
)

type SchDitem struct {
	Id string //PdplinkHash
	// PdplinkHash string
	TopUrl     string
	TargetPath string // URL without scheme info
	HtmlHash   string
	HubUrls    string
	IgroupCode string
	IgroupName string
	Categories string
	Emails     string
	ItemName   string
	Tags       string
	Currency   string
	SalesPrice string

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

type SchDitemRow struct {
	ID     string   `json:"_id"`
	Index  string   `json:"_index"`
	Score  float64  `json:"_score"`
	Source SchDitem `json:"_source"`
	Type   string   `json:"_type"`
}

func (t *SchDitem) TableName() string { // 반드시 있어야 table name을 가져올 수 있다.
	return e.TableName(*t)
}

func (t *SchDitem) RecordName() string { //필수
	return "Id of " + t.Id + " in " + e.TableName(*t) + " "
}

func (t *SchDitem) GetaRow(y *elasticsearch.Client) error {

	res, err := y.Get(e.TableName(*t), t.Id)
	defer res.Body.Close()

	if err != nil {
		return e.LogErr("qrlhbdft", "Error of index: "+t.TableName()+": ", err)
	} else {
		if res.IsError() { //NotFound
			return errors.New(e.RecNotFound("qrhbdfke4k", "in "+t.TableName()+" of Id: "+t.Id))
		}
	}

	var resJs SchDitemRow
	if err := json.NewDecoder(res.Body).Decode(&resJs); err != nil {
		return e.LogErr("qrlhbdfk", "Error decoding "+t.TableName()+": ", err)
	}

	*t = SchDitem(resJs.Source)
	// *t = resJs.Source
	return nil
}

// func (t *SchDitem) AddaRow(y *elasticsearch.Client) error {

// 	data, _ := json.Marshal(t)
// 	res, err := y.Index(
// 		e.TableName(*t),
// 		bytes.NewReader(data),
// 		y.Index.WithDocumentID((t.Id)))
// 	defer res.Body.Close()

// 	if err != nil {
// 		return errors.New(e.RecAddErr("jfhgytg", t.RecordName()+err.Error()))
// 	}
// 	return nil

// }

// func (t *SchDitem) DelaRow(y *elasticsearch.Client) error {

// 	res, err := y.Delete(e.TableName(*t), t.Id)
// 	defer res.Body.Close()

// 	if err != nil {
// 		return errors.New(e.RecAddErr("jfhgy4tg", t.RecordName()+err.Error()))
// 	}
// 	return nil
// }
