package config

import (
	"encoding/json"
	"log"
	"net/http"
	"newshub/allstruct"
	"os"
	"path/filepath"
	"slices"
)

var (
	VersionName string = "v0.0.8"
	VersionCode int    = 2026062000

	RunDir      string
	DataBaseDir string = "/database"
	AnalysisDir string = "/analysis"

	Ok bool

	UserName     string = "rexxrt"
	UserPassword string = "PleaseInputText"
	Logined      []string

	Save2DBTime int64 = 60 * 15

	Models map[string]allstruct.ModelInfo = make(map[string]allstruct.ModelInfo)

	SourceZhihu    string = "https://orz.ai/api/v1/dailynews/multi?platforms=zhihu"
	SourceBilibili string = "https://orz.ai/api/v1/dailynews/multi?platforms=bilibili"
	SourceCoolapk  string = "https://newsnow.danhua.ddns-ip.net/api/s?id=coolapk"
	SourceNewsNow  string = "https://newsnow.busiyi.world/api/s?id="
	SourceNewsNow2 string = "https://newsnow.busiyi.world/api/s/entire"
)

func Init() {
	var configpath string = filepath.Join(RunDir, "/config.json")
	if data, err := os.ReadFile(configpath); err == nil {
		fconfig := make(map[string]any)
		if json.Unmarshal(data, &fconfig) == nil {
			if v, ok := fconfig["save2dbtime"].(float64); ok {
				Save2DBTime = int64(v)
			}
			if v, ok := fconfig["sourcenewsnow"].(string); ok {
				SourceNewsNow = v
			}
			if v, ok := fconfig["sourcenewsnow2"].(string); ok {
				SourceNewsNow2 = v
			}
			if v, ok := fconfig["username"].(string); ok {
				UserName = v
			}
			if v, ok := fconfig["userpassword"].(string); ok {
				UserPassword = v
			}
			if v, ok := fconfig["models"].(map[string]any); ok {
				log.Println("品鉴中")
				for k, v2 := range v {
					v2 := v2.(map[string]any)
					log.Println(k, v2)
					var model allstruct.ModelInfo
					model.Model = v2["Model"].(string)
					model.Key = v2["Key"].(string)
					model.Url = v2["Url"].(string)
					Models[k] = model
				}
			} else {
				log.Println("无法解析 models")
			}
		} else {
			log.Println("无法解析配置 config.json")
		}
	} else {
		log.Println("无法读入设置 config.json，已进行初始化")
		SaveSettings()
	}
}

func SaveSettings() {
	var configpath string = filepath.Join(RunDir, "/config.json")
	fconfig := map[string]any{
		"save2dbtime":    Save2DBTime,
		"sourcenewsnow":  SourceNewsNow,
		"sourcenewsnow2": SourceNewsNow2,
		"username":       UserName,
		"userpassword":   UserPassword,
		"models":         Models,
	}
	if sdata, err := json.MarshalIndent(fconfig, "", "    "); err == nil {
		os.WriteFile(configpath, sdata, 0755)
	}
}

func CheckLogin(r *http.Request) bool {
	if cookie, err := r.Cookie("Authorization"); err == nil {
		if slices.Contains(Logined, cookie.Value) {
			return true
		}
	}
	return false
}
