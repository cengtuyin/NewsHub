package news

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"newshub/config"
	"newshub/database"
	"time"
)

var (
	lastinsert int64
)

func Init() {

}

func GetNewsJson(method string, sourceapi string) (map[string]any, error) {
	var (
		client http.Client
		rdata  map[string]any
	)
	req, _ := http.NewRequest("GET", sourceapi, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36 Edg/149.0.0.0")
	if r, err := client.Do(req); err == nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			if json.Unmarshal(data, &rdata) == nil {
				return rdata, nil
			} else {
				return nil, fmt.Errorf("GetNewsJson Unmarshal (%s) : %s", sourceapi, string(data))
			}
		} else {
			return nil, fmt.Errorf("GetNewsJson ReadAll (%s) : %s", sourceapi, err)
		}
	} else {
		return nil, fmt.Errorf("GetNewsJson Get (%s) : %s", sourceapi, err)
	}
}
func GetNewsJson2(method string, sourceapi string) ([]any, error) {
	var (
		client http.Client
		rdata  []any
	)
	req, _ := http.NewRequest(method, sourceapi, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36 Edg/149.0.0.0")
	if r, err := client.Do(req); err == nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			if json.Unmarshal(data, &rdata) == nil {
				return rdata, nil
			} else {
				return nil, fmt.Errorf("GetNewsJson Unmarshal (%s) : %s", sourceapi, string(data))
			}
		} else {
			return nil, fmt.Errorf("GetNewsJson ReadAll (%s) : %s", sourceapi, err)
		}
	} else {
		return nil, fmt.Errorf("GetNewsJson Get (%s) : %s", sourceapi, err)
	}
}
func GetNewsJsonPOST(sourceapi string) ([]any, error) {
	var (
		client http.Client
		rdata  []any
	)
	req, _ := http.NewRequest("POST", sourceapi, bytes.NewBuffer([]byte(`{"sources":["36kr-renqi","baidu","bilibili-hot-search","chongbuluo-hot","cls-hot","coolapk","douyin","freebuf","github-trending-today","hackernews","hupu","ifeng","iqiyi-hot-ranklist","juejin","nowcoder","producthunt","qqvideo-tv-hotsearch","steam","tencent-hot","thepaper","tieba","toutiao","wallstreetcn-hot","weibo","zhihu"]}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36 Edg/149.0.0.0")
	if r, err := client.Do(req); err == nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			if json.Unmarshal(data, &rdata) == nil {
				return rdata, nil
			} else {
				return nil, fmt.Errorf("GetNewsJson Unmarshal (%s) : %s", sourceapi, string(data))
			}
		} else {
			return nil, fmt.Errorf("GetNewsJson ReadAll (%s) : %s", sourceapi, err)
		}
	} else {
		return nil, fmt.Errorf("GetNewsJson Get (%s) : %s", sourceapi, err)
	}
}

/*
	{
		"title": "黑龙江 5 名高中生自主研发的临近空间火箭发射成功，这是什么水平？给高中教育带来哪些启示？",
		"url": "https://www.zhihu.com/question/2048037405521981848",
		"content": "四问黑龙江共青团。前排：评论不是我删的。第一问，高能燃料哪里来？黑龙江共青团的报道说，火箭用固体燃料是高中生自学教材、自研出来的： [图片] 火箭高能燃料，尤其是报道里的固体燃料，在化学本质和物理特性上完全等同于爆炸物。制造固体火箭燃料必不可少的氧化剂（比如高氯酸铵、硝酸钾）和还原剂（比如…",
		"source": "zhihu",
		"publish_time": "2026-06-11 10:46:28"
	}
*/
func GetBilibiliNews() ([]any, error) {
	if data, err := GetNewsJson("GET", config.SourceBilibili); err == nil {
		if data["status"].(string) == "200" {
			return data["data"].(map[string]any)["bilibili"].([]any), nil
		} else {
			return nil, fmt.Errorf("Source Error")
		}
	} else {
		return nil, err
	}
}

/*
	{
		"title": "末世先杀圣母",
		"url": "https://www.bilibili.com/video/BV1keEu61EfE",
		"content": "-",
		"source": "bilibili",
		"publish_time": "2026-06-11 10:46:34"
	}
*/
func GetZhihuNews() ([]any, error) {
	if data, err := GetNewsJson("GET", config.SourceZhihu); err == nil {
		if data["status"].(string) == "200" {
			return data["data"].(map[string]any)["zhihu"].([]any), nil
		} else {
			return nil, fmt.Errorf("Source Error")
		}
	} else {
		return nil, err
	}
}

/*
	{
		"id": 72263472,
		"title": "[捂脸][捂脸][捂脸]高考完当天就被小米偷家了，果然是针对高考生的手机。",
		"url": "",
		"extra": {
			"info": "134.6万热度 700讨论"
		}
	}
*/
func GetCoolapkNews() ([]any, error) {
	if data, err := GetNewsJson("GET", config.SourceCoolapk); err == nil {
		if data["status"].(string) == "success" {
			return data["items"].([]any), nil
		} else {
			return nil, fmt.Errorf("Source Error")
		}
	} else {
		return nil, err
	}
}

func GetNewsNow(tag string) ([]any, error) {
	if data, err := GetNewsJson("GET", config.SourceNewsNow+tag); err == nil {
		if data["status"].(string) == "success" || data["status"].(string) == "cache" {
			return data["items"].([]any), nil
		} else {
			return nil, fmt.Errorf("Source Error")
		}
	} else {
		return nil, err
	}
}

func GetNewsNow2() ([]any, error) {
	if data, err := GetNewsJsonPOST(config.SourceNewsNow2); err == nil {
		sdata := make(map[string][]byte)
		for _, t := range data {
			news := t.(map[string]any)
			if news["id"] == "zhihu" || news["id"] == "coolapk" || news["id"] == "juejin" || news["id"] == "weibo" || news["id"] == "toutiao" || news["id"] == "thepaper" || news["id"] == "hackernews" || news["id"] == "cls-hot" || news["id"] == "freebuf" {
				if news["status"] == "cache" || news["status"] == "success" {
					if ssdata, err := json.Marshal(news["items"]); err == nil {
						sdata[news["id"].(string)] = ssdata
					}
				}
			}
		}
		nowtime := time.Now().Unix()
		if nowtime-lastinsert > config.Save2DBTime {
			lastinsert = nowtime
			go database.Newsdb.Exec(`INSERT INTO news (zhihu,coolapk,juejin,weibo,toutiao,thepaper,hackernews,clshot,freebuf) VALUES (?,?,?,?,?,?,?,?,?)`, sdata["zhihu"], sdata["coolapk"], sdata["juejin"], sdata["weibo"], sdata["toutiao"], sdata["thepaper"], sdata["hackernews"], sdata["cls-hot"], sdata["freebuf"])
		}
		return data, nil
	} else {
		return nil, err
	}
}

func GetAllNews() {

}
