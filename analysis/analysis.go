package analysis

import (
	"encoding/json"
	"fmt"
	"log"
	"newshub/config"
	"newshub/database"
	"newshub/news"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	WordsCloud   []string
	newsnow2tmp  []any
	newsnow2tmp2 []string
	lastinsert   []int64 = []int64{0, 0}
)

func Init() {
	config.AnalysisDir = filepath.Join(config.RunDir, config.AnalysisDir)
	os.MkdirAll(config.AnalysisDir, 0755)
	if data, err := os.ReadFile(filepath.Join(config.AnalysisDir, "/wordscloud.txt")); err == nil {
		var banword []string
		if data, err := os.ReadFile(filepath.Join(config.AnalysisDir, "/wordscloud_ban.txt")); err == nil {
			banword = strings.Split(string(data), "\n")
		}
		WordsCloud = strings.Split(string(data), "\n")
		banMap := make(map[string]bool)
		for _, w := range banword {
			banMap[strings.TrimSpace(w)] = true
		}
		for i := len(WordsCloud) - 1; i >= 0; i-- {
			word := strings.TrimSpace(WordsCloud[i])
			if banMap[word] {
				WordsCloud = append(WordsCloud[:i], WordsCloud[i+1:]...)
				continue
			}
			if utf8.RuneCountInString(word) < 2 {
				WordsCloud = append(WordsCloud[:i], WordsCloud[i+1:]...)
			}
		}
		sort.SliceStable(WordsCloud, func(i, j int) bool {
			return utf8.RuneCountInString(WordsCloud[i]) > utf8.RuneCountInString(WordsCloud[j])
		})
		log.Printf("词云载入词汇 %d 条\n", len(WordsCloud))
	} else {
		log.Println("无法载入词云库 /analysis/wordscloud.txt")
	}
}

func Words(length int) (map[string]int, error) {
	var (
		rdata  map[string]int = make(map[string]int)
		strtmp string
	)
	if data, err := news.GetNewsNow2(); err == nil {
		newsnow2tmp = data
		for _, t := range data {
			news := t.(map[string]any)
			if news["status"].(string) == "success" || news["status"].(string) == "cache" {
				for _, thenew2 := range news["items"].([]any) {
					thenew := thenew2.(map[string]any)
					if tm, ok := thenew["title"].(string); ok {
						strtmp += tm
					}
					if tm, ok := thenew["extra"].(map[string]any); ok {
						if tm, ok := tm["hover"].(string); ok {
							strtmp += " / "
							strtmp += tm
						}
					}
				}
			}
		}
		for _, v := range WordsCloud {
			if v == "" || unicode.IsPunct(rune(v[0])) || unicode.IsSymbol(rune(v[0])) {
				continue
			}
			i := strings.Count(strtmp, v)
			if i > length {
				rdata[v] = i
			}
		}
		newsnow2tmp2 = []string{}
		for k := range rdata {
			newsnow2tmp2 = append(newsnow2tmp2, k)
		}
		nowtime := time.Now().Unix()
		if nowtime-lastinsert[0] > config.Save2DBTime {
			lastinsert[0] = nowtime
			if sdata, err := json.Marshal(rdata); err == nil {
				go database.WordsClouddb.Exec(`INSERT INTO wordscloud (data) VALUES (?)`, sdata)
			}
		}
		return rdata, nil
	} else {
		return nil, err
	}
}

func WordFindNews(word string) ([]map[string]any, error) {
	var (
		rdata []map[string]any
	)
	if len(newsnow2tmp) == 0 {
		if _, err := Words(10); err != nil {
			return nil, fmt.Errorf("Unable get newsnow2")
		}
	}
	for _, t := range newsnow2tmp {
		news := t.(map[string]any)
		if news["status"].(string) == "success" || news["status"].(string) == "cache" {
			for _, thenew2 := range news["items"].([]any) {
				thenew := thenew2.(map[string]any)
				if tm, ok := thenew["title"].(string); ok {
					if strings.Count(tm, word) > 0 {
						rdata = append(rdata, thenew)
						continue
					}
				}
				if tm, ok := thenew["extra"].(map[string]any); ok {
					if tm, ok := tm["hover"].(string); ok {
						if strings.Count(tm, word) > 0 {
							rdata = append(rdata, thenew)
						}
					}
				}
			}
		}
	}
	return rdata, nil
}

func WordFindNewss() ([]any, map[string][]int, error) {
	var fromnews []any
	var words map[string][]int = make(map[string][]int)
	if len(newsnow2tmp2) == 0 {
		if _, err := Words(10); err != nil {
			return nil, nil, fmt.Errorf("unable get newsnow2: %v", err)
		}
	}
	for _, word := range newsnow2tmp2 {
		if list, err := WordFindNews(word); err == nil {
			for _, v := range list {
				in := false
				for id, v2 := range fromnews {
					if reflect.DeepEqual(v, v2) {
						in = true
						words[word] = append(words[word], id)
					}
				}
				if !in {
					fromnews = append(fromnews, v)
					words[word] = append(words[word], len(fromnews)-1)
				}
			}
		}
	}
	return fromnews, words, nil
}

func WordFindNews2(length int) ([]any, map[string][]int, error) {
	var thenews []any
	rdata := make(map[string][]int)
	newsCache := make(map[string]bool)
	if len(newsnow2tmp2) == 0 {
		if _, err := Words(10); err != nil {
			return nil, nil, fmt.Errorf("unable get newsnow2: %v", err)
		}
	}
	for i, word := range newsnow2tmp2 {
		tmp, err := WordFindNews(word)
		if err != nil {
			continue
		}
		for _, thenew := range tmp {
			title, _ := thenew["title"].(string)
			var hover string
			if extra, ok := thenew["extra"].(map[string]any); ok {
				hover, _ = extra["hover"].(string)
			}
			var matchedWords []string
			for j := i; j < len(newsnow2tmp2); j++ {
				word2 := newsnow2tmp2[j]
				if strings.Contains(title, word2) || strings.Contains(hover, word2) {
					matchedWords = append(matchedWords, word2)
				}
			}
			if len(matchedWords) < 2 {
				continue
			}
			//sort.Strings(matchedWords)
			sort.SliceStable(matchedWords, func(i, j int) bool {
				return utf8.RuneCountInString(matchedWords[i]) > utf8.RuneCountInString(matchedWords[j])
			})
			for _, v := range matchedWords {
				in := -1
				for i, v2 := range matchedWords {
					if v != v2 && strings.Count(v, v2) != 0 {
						in = i
					}
				}
				if in != -1 {
					matchedWords = append(matchedWords[:in], matchedWords[in+1:]...)
				}
			}
			key := strings.Join(matchedWords, "|")
			cacheKey := key + "|" + title + "|" + hover
			if newsCache[cacheKey] {
				continue
			}
			newsCache[cacheKey] = true
			thenews = append(thenews, thenew)
			rdata[key] = append(rdata[key], len(thenews)-1)
		}
	}

	for k, keys := range rdata {
		if len(keys) < length {
			delete(rdata, k)
		}
	}

	nowtime := time.Now().Unix()
	if nowtime-lastinsert[1] > config.Save2DBTime {
		lastinsert[1] = nowtime
		sdata := make(map[string][]any)
		for k, v := range rdata {
			var vs []any
			for _, v := range v {
				vs = append(vs, thenews[v])
			}
			sdata[k] = vs
		}
		if sdata, err := json.Marshal(sdata); err == nil {
			go database.LinkWords.Exec(`INSERT INTO linkwords (data) VALUES (?)`, sdata)
		}
	}

	var thenews2 []any
	for rk, keys := range rdata {
		for rkk, v := range keys {
			in := -1
			for id, t := range thenews2 {
				if reflect.DeepEqual(t, thenews[v]) {
					in = id
					break
				}
			}
			if in == -1 {
				thenews2 = append(thenews2, thenews[v])
				rdata[rk][rkk] = len(thenews2) - 1
			} else {
				rdata[rk][rkk] = in
			}
		}
	}
	return thenews2, rdata, nil
}
