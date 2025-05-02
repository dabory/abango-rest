package abango

import (
	"encoding/json"
	"fmt"
	"os"

	e "github.com/dabory/abango-rest/etc"

	_ "github.com/go-sql-driver/mysql"
)

func GetXConfig(params ...string) error { // Kafka, gRpc, REST 통합 업그레이드

	e.AokLog("Abango Gets XConfig !")
	conf := "conf/"
	if len(params) != 0 {
		conf = params[0] + conf
	}

	RunFilename := conf + "config_select.json"

	run := struct {
		ConfSelect  string
		ConfPostFix string
	}{}

	if file, err := os.Open(RunFilename); err != nil {
		e.LogFatal("WERQRRQERQW", RunFilename+"  File NOT exist", err)
		return err
	} else {
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&run); err != nil {
			e.LogFatal("ERTFDFDAFA", RunFilename+"Not Decoded", err)
			return err
		}
	}

	XConfig = make(map[string]string) // Just like malloc
	config := []Param{}

	// var varMap []map[string]interface{}
	filename := conf + run.ConfSelect + run.ConfPostFix
	if file, err := os.Open(filename); err != nil {
		e.LogFatal("WERCASDFAWEF", filename+"  File NOT exist", err)
		return err
	} else {
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&config); err == nil {
			for _, p := range config {
				XConfig[p.Key] = p.Value
			}
		} else {
			e.LogFatal("QWERCQQGE", filename+"Not Decoded", err)
		}
	}

	if XConfig["RestOn"] == "Yes" || XConfig["ApiType"] == "Rest" {
		fmt.Println("==" + "Config file prefix: " + run.ConfSelect + "== REST Connection: " + XConfig["RestConnect"] + "==")
	}
	return nil
}
