package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"newshub/analysis"
	"newshub/config"
	"newshub/database"
	"newshub/handle"
	"newshub/model"
	"newshub/news"
	"newshub/tasks"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// 定义数据结构
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 模拟数据库
var users = []User{
	{ID: "1", Name: "Alice", Age: 25},
	{ID: "2", Name: "Bob", Age: 30},
}

func main() {
	log.Println("NewsHub Start!")
	log.Println("Made By Rexxrt.")

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	config.RunDir = filepath.Dir(exe)

	log.Println("载入配置文件")
	config.Init()
	log.Println("初始化数据库")
	database.Init()
	log.Println("初始化新闻")
	news.Init()
	log.Println("初始化请求")
	handle.Init()
	log.Println("初始化分析")
	analysis.Init()
	log.Println("初始化模型")
	model.Init()
	log.Println("初始化任务")
	tasks.Init()
	log.Println("启动HTTP服务")

	// 创建多路复用器
	r := http.NewServeMux()

	// 注册路由
	r.HandleFunc("/api/bilibili/news", handle.GetBilibiliNews)
	r.HandleFunc("/api/zhihu/news", handle.GetZhihuNews)
	r.HandleFunc("/api/coolapk/news", handle.GetCoolapkNews)
	r.HandleFunc("/api/newsnow", handle.GetNewsNow)
	r.HandleFunc("/api/newsnow2", handle.GetNewsNow2)
	r.HandleFunc("/api/wordscloud", handle.WordsCloud)
	r.HandleFunc("/api/wordsearch", handle.WordFindNews)
	r.HandleFunc("/api/wordsearch2", handle.WordFindNewss)
	r.HandleFunc("/api/wordsin", handle.WordFindNews2)
	r.HandleFunc("/api/login", handle.Login)
	r.HandleFunc("/api/settings", handle.UpdateSettings)
	r.HandleFunc("/api/chat", handle.Chat)
	r.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		data := map[string]any{
			"application": "NewsHub",
			"versionName": config.VersionName,
			"versionCode": config.VersionCode,
		}
		rdata, _ := json.Marshal(data)
		w.Write(rdata)
	})

	// 静态文件服务
	r.Handle("/", http.FileServer(http.Dir(filepath.Join(config.RunDir, "/static"))))

	// 中间件链
	handler := loggingMiddleware(r)
	handler = corsMiddleware(handler)
	handler = authorizationMiddleware(handler)

	// 启动服务器
	server := &http.Server{
		Addr:    ":51516",
		Handler: handler,
	}

	log.Println("Server starting on http://localhost:51516")
	config.Ok = true
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, delIPPort(GetRealIP(r)))
		next.ServeHTTP(w, r)
	})
}

// CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Authorization 中间件
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logined := config.CheckLogin(r)
		path := strings.ToLower(r.URL.Path)
		if slices.Contains([]string{"/api/settings", "/api/chat"}, path) && !logined {
			logintips(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logintips(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(r.URL.Path)
	if strings.HasPrefix(path, "/api") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		if data, err := json.Marshal(map[string]any{
			"success": false,
			"message": "需要登录后才可操作",
		}); err == nil {
			if _, err := w.Write(data); err != nil {
				log.Println("Words Write Error")
			}
		}
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		http.ServeFile(w, r, filepath.Join(config.RunDir, "/static/403.html"))
	}
}

// 获取真实客户端 IP
func GetRealIP(r *http.Request) string {
	for _, header := range []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai
		"X-Client-IP",
		"Forwarded"} {
		if xff := r.Header.Get(header); xff != "" {
			ips := strings.Split(xff, ",")
			if len(ips) > 0 {
				ip := strings.TrimSpace(ips[0])
				if isValidIP(ip) {
					return ip
				}
			}
		}
	}

	// 检查 X-Real-IP（Nginx 常用）
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if isValidIP(xri) {
			return xri
		}
	}

	// 检查 Forwarded（RFC 7239）
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		// 解析 Forwarded 头
		for _, part := range strings.Split(forwarded, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "for=") {
				ip := strings.TrimPrefix(part, "for=")
				ip = strings.Trim(ip, "\"")
				if isValidIP(ip) {
					return ip
				}
			}
		}
	}

	// 最后从 RemoteAddr 获取
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func delIPPort(ip string) string {
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	// 去除 IPv6 的方括号
	ip = strings.TrimPrefix(ip, "[")
	ip = strings.TrimSuffix(ip, "]")
	return ip
}

func isValidIP(ip string) bool {
	return net.ParseIP(delIPPort(ip)) != nil
}
