package handle

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"newshub/allstruct"
	"newshub/analysis"
	"newshub/config"
	"newshub/model"
	"newshub/news"
	"strconv"
	"sync"
)

var (
	LastWordsCloudData []byte
	WordsCloudAction   bool
	LastNewsNow2Data   []byte
	NewsNow2Action     bool
	LastNewNowData     map[string][]byte
	NewsNowAction      map[string]bool
)

func Init() {
	LastNewNowData = make(map[string][]byte)
	NewsNowAction = make(map[string]bool)
}

func GetBilibiliNews(w http.ResponseWriter, r *http.Request) {
	if data, err := news.GetBilibiliNews(); err == nil {
		if data, err := json.Marshal(data); err == nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			if _, err := w.Write(data); err != nil {
				log.Println("GetBilibiliNews Write Error")
			}
		}
	}
}

func GetZhihuNews(w http.ResponseWriter, r *http.Request) {
	if data, err := news.GetZhihuNews(); err == nil {
		if data, err := json.Marshal(data); err == nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			if _, err := w.Write(data); err != nil {
				log.Println("GetZhihuNews Write Error")
			}
		}
	}
}

func GetCoolapkNews(w http.ResponseWriter, r *http.Request) {
	if data, err := news.GetCoolapkNews(); err == nil {
		if data, err := json.Marshal(data); err == nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			if _, err := w.Write(data); err != nil {
				log.Println("GetCoolapkNews Write Error")
			}
		}
	}
}

func GetNewsNow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		returnMessage(w, false, "需要GET参数[id]")
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if len(LastNewNowData[id]) != 0 {
		w.Write(LastNewNowData[id])
	}
	if NewsNowAction[id] {
		if len(LastNewNowData[id]) == 0 {
			returnMessage(w, false, "任务正在进行中")
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	NewsNowAction[id] = true
	go func() {
		defer func() {
			NewsNowAction[id] = false
			wg.Done()
		}()
		if data, err := news.GetNewsNow(id); err == nil {
			if data, err := json.Marshal(data); err == nil {
				if len(LastNewNowData[id]) == 0 {
					if _, err := w.Write(data); err != nil {
						log.Println("Words Write Error")
					}
				}
				LastNewNowData[id] = data
			} else {
				returnMessage(w, false, fmt.Sprintf("%s", err))
				fmt.Fprint(w, err)
			}
		} else {
			returnMessage(w, false, fmt.Sprintf("%s", err))
			fmt.Fprint(w, err)
		}
	}()
	if len(LastNewNowData[id]) == 0 {
		wg.Wait()
	}
}

func GetNewsNow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.URL.Query().Get("now") != "" {
		LastNewsNow2Data = []byte{}
	}
	if len(LastNewsNow2Data) != 0 {
		w.Write(LastNewsNow2Data)
	}
	if NewsNow2Action {
		if len(LastNewsNow2Data) == 0 {
			returnMessage(w, false, "任务正在进行中")
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	NewsNow2Action = true
	go func() {
		defer func() {
			NewsNow2Action = false
			wg.Done()
		}()
		if data, err := news.GetNewsNow2(); err == nil {
			if data, err := json.Marshal(data); err == nil {
				if len(LastNewsNow2Data) == 0 {
					if _, err := w.Write(data); err != nil {
						log.Println("Words Write Error")
					}
				}
				LastNewsNow2Data = data
			} else {
				returnMessage(w, false, fmt.Sprintf("%s", err))
				fmt.Fprint(w, err)
			}
		} else {
			returnMessage(w, false, fmt.Sprintf("%s", err))
			fmt.Fprint(w, err)
		}
	}()
	if len(LastNewsNow2Data) == 0 {
		wg.Wait()
	}
}

func WordsCloud(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.URL.Query().Get("length") == "" {
		returnMessage(w, false, "需要GET参数[length]")
		return
	}
	if r.URL.Query().Get("now") != "" {
		LastWordsCloudData = []byte{}
	}
	if len(LastWordsCloudData) != 0 {
		w.Write(LastWordsCloudData)
	}
	if WordsCloudAction {
		if len(LastWordsCloudData) == 0 {
			returnMessage(w, false, "任务正在进行中")
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	WordsCloudAction = true
	go func() {
		defer func() {
			WordsCloudAction = false
			wg.Done()
		}()
		if target, ok := strconv.Atoi(r.URL.Query().Get("length")); ok == nil {
			var rdata map[string]any = make(map[string]any)
			if data, err := analysis.Words(target); err == nil {
				rdata["success"] = true
				rdata["length"] = len(data)
				rdata["words"] = data
				if data, err := json.Marshal(rdata); err == nil {
					if len(LastWordsCloudData) == 0 {
						if _, err := w.Write(data); err != nil {
							log.Println("Words Write Error")
						}
					}
					LastWordsCloudData = data
				} else {
					returnMessage(w, false, fmt.Sprintf("%s", err))
				}
			} else {
				returnMessage(w, false, fmt.Sprintf("%s", err))
			}
		} else {
			returnMessage(w, false, "GET参数[length]需要为整数")
		}
	}()
	if len(LastWordsCloudData) == 0 {
		wg.Wait()
	}
}

func WordFindNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.URL.Query().Get("word") == "" {
		returnMessage(w, false, "需要GET参数[word]")
		return
	}
	var rdata map[string]any = make(map[string]any)
	if data, err := analysis.WordFindNews(r.URL.Query().Get("word")); err == nil {
		rdata["success"] = true
		rdata["length"] = len(data)
		rdata["news"] = data
		if data, err := json.Marshal(rdata); err == nil {
			if _, err := w.Write(data); err != nil {
				log.Println("Words Write Error")
			}
		} else {
			returnMessage(w, false, fmt.Sprintf("%s", err))
		}
	} else {
		returnMessage(w, false, fmt.Sprintf("%s", err))
	}
}

func WordFindNewss(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	var rdata map[string]any = make(map[string]any)
	if adata, data, err := analysis.WordFindNewss(); err == nil {
		rdata["success"] = true
		rdata["news_length"] = len(adata)
		rdata["words_length"] = len(data)
		rdata["news"] = adata
		rdata["words"] = data
		if data, err := json.Marshal(rdata); err == nil {
			if _, err := w.Write(data); err != nil {
				log.Println("Words Write Error")
			}
		} else {
			returnMessage(w, false, fmt.Sprintf("%s", err))
		}
	} else {
		returnMessage(w, false, fmt.Sprintf("%s", err))
	}
}

func WordFindNews2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.URL.Query().Get("length") == "" {
		returnMessage(w, false, "需要GET参数[length]")
		return
	}

	if target, ok := strconv.Atoi(r.URL.Query().Get("length")); ok == nil {
		var rdata map[string]any = make(map[string]any)
		if adata, data, err := analysis.WordFindNews2(target); err == nil {
			rdata["success"] = true
			rdata["news_length"] = len(adata)
			rdata["words_length"] = len(data)
			rdata["news"] = adata
			rdata["words"] = data
			if data, err := json.Marshal(rdata); err == nil {
				if _, err := w.Write(data); err != nil {
					log.Println("Words Write Error")
				}
			} else {
				returnMessage(w, false, fmt.Sprintf("%s", err))
			}
		} else {
			returnMessage(w, false, fmt.Sprintf("%s", err))
		}
	} else {
		returnMessage(w, false, "GET参数[length]需要为整数")
	}
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		keys := r.URL.Query().Get("keys")
		if keys == "" {
			returnMessage(w, false, "需要GET参数[keys]")
			return
		}
		var klist []any
		if err := json.Unmarshal([]byte(keys), &klist); err != nil {
			returnMessage(w, false, "提交的数据错误")
			return
		}
		defer r.Body.Close()
		rdata := make(map[string]any)
		var i int
		for i = 0; i < len(klist); i++ {
			key := klist[i]
			if keys, ok := key.([]string); ok {
				switch keys[0] {
				case "Models":
					rdata["Models"] = []map[string]string{}
					for _, v := range keys {
						rdatas := make(map[string]string)
						if v != keys[0] {
							rdatas["Url"] = config.Models[v].Url
							rdatas["Key"] = config.Models[v].Key
							rdatas["Model"] = config.Models[v].Model
						}
						rdata["Models"] = append(rdata["Models"].([]map[string]string), rdatas)
					}
				default:
					returnMessage(w, false, "不受支持的[key]")
					return
				}
			} else if key, ok := key.(string); ok && key != "" {
				switch key {
				case "Save2DBTime":
					rdata[key] = config.Save2DBTime
				case "SourceNewsNow":
					rdata[key] = config.SourceNewsNow
				case "SourceNewsNow2":
					rdata[key] = config.SourceNewsNow2
				default:
					returnMessage(w, false, "不受支持的[key]")
					return
				}
			}
		}
		if data, err := json.Marshal(map[string]any{
			"success": true,
			"data":    rdata,
		}); err == nil {
			if _, err := w.Write(data); err != nil {
				log.Println("Words Write Error")
			}
		}
	case http.MethodPost:
		var data map[string]any
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			returnMessage(w, false, "提交的数据错误")
			return
		}
		defer r.Body.Close()
		klist, ok := data["keys"].([]any)
		if !ok {
			returnMessage(w, false, "需要POST&JSON参数[keys]")
			return
		}
		vlist, ok := data["values"].([]any)
		if !ok {
			returnMessage(w, false, "需要POST&JSON参数[values]")
			return
		}
		var i int
		for i = 0; i < len(klist); i++ {
			key := klist[i]
			value := vlist[i]
			if keys, ok := key.([]string); ok {
				switch keys[0] {
				case "Models":
					values := value.([]map[string]string)
					for i, v := range keys {
						var modelinfo allstruct.ModelInfo
						if v != keys[0] {
							if vvv, ok := values[i]["Url"]; ok {
								modelinfo.Url = vvv
							}
							if vvv, ok := values[i]["Key"]; ok {
								modelinfo.Key = vvv
							}
							if vvv, ok := values[i]["Model"]; ok {
								modelinfo.Model = vvv
							}
						}
						config.Models[values[i]["Name"]] = modelinfo
					}
					config.SaveSettings()
				default:
					returnMessage(w, false, "不受支持的[key]")
					return
				}
			} else if key, ok := key.(string); ok && key != "" {
				switch key {
				case "Save2DBTime":
					if target, ok := value.(float64); ok {
						config.Save2DBTime = int64(target)
						config.SaveSettings()
					} else {
						returnMessage(w, false, key+"参数[value]需要为整数")
						return
					}
				case "SourceNewsNow":
					config.SourceNewsNow = value.(string)
					config.SaveSettings()
				case "SourceNewsNow2":
					config.SourceNewsNow2 = value.(string)
					config.SaveSettings()
				default:
					returnMessage(w, false, "不受支持的[key]:"+key)
					return
				}
			}
		}
		returnMessage(w, true, "修改成功")
	}
}

func returnData(w http.ResponseWriter, status bool, data any) {
	if data, err := json.Marshal(map[string]any{
		"success": status,
		"data":    data,
	}); err == nil {
		if _, err := w.Write(data); err != nil {
			log.Println("Words Write Error")
		}
	}
}

func returnMessage(w http.ResponseWriter, status bool, message string) {
	if data, err := json.Marshal(map[string]any{
		"success": status,
		"message": message,
	}); err == nil {
		if _, err := w.Write(data); err != nil {
			log.Println("Words Write Error")
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	var (
		user     string = r.URL.Query().Get("user")
		password string = r.URL.Query().Get("password")
	)
	if user == "" {
		returnMessage(w, false, "需要GET参数[user]")
		return
	}
	if password == "" {
		returnMessage(w, false, "需要GET参数[password]")
		return
	}
	if user == config.UserName {
		if password == config.UserPassword {
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|:.<>?/~"
			b := make([]byte, 32)
			for i := range b {
				num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
				b[i] = charset[num.Int64()]
			}
			config.Logined = append(config.Logined, string(b))
			cookie := &http.Cookie{
				Name:   "Authorization",
				Value:  string(b),
				Path:   "/",
				MaxAge: 60 * 60 * 3,
			}
			http.SetCookie(w, cookie)
			returnMessage(w, true, "登录成功")
			return
		}
	}
	returnMessage(w, false, "用户名或密码错误")
}

func Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		smodel := r.URL.Query().Get("model")
		if smodel == "" {
			returnMessage(w, false, "需要GET参数[model]")
			return
		}
		message := r.URL.Query().Get("message")
		if message == "" {
			returnMessage(w, false, "需要GET参数[message]或POST&JSON（ChatAPI风格）")
			return
		}
		rdata, err := model.Chat(smodel, []map[string]string{
			{
				"role":    "user",
				"content": message,
			},
		})
		if err != nil {
			returnMessage(w, false, err.Error())
			return
		}
		returnData(w, true, rdata)
	} else {
		smodel := r.URL.Query().Get("model")
		if smodel == "" {
			returnMessage(w, false, "需要GET参数[model]")
			return
		}
		var data []map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			returnMessage(w, false, "提交的数据错误")
			return
		}
		defer r.Body.Close()
		rdata, err := model.Chat(smodel, data)
		if err != nil {
			returnMessage(w, false, err.Error())
			return
		}
		returnData(w, true, rdata)
	}
}
