package database

import (
	"database/sql"
	"log"
	"newshub/config"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var (
	Newsdb       *sql.DB
	WordsClouddb *sql.DB
	LinkWords    *sql.DB
)

func Init() {
	config.DataBaseDir = filepath.Join(config.RunDir, config.DataBaseDir)
	os.MkdirAll(config.DataBaseDir, 0755)

	_init_newsdb()
	_init_wordsclouddb()
	_init_linkwords()
}

func _init_newsdb() {
	var err error
	Newsdb, err = sql.Open("sqlite", filepath.Join(config.DataBaseDir, "/news.db"))
	if err != nil {
		log.Println("无法打开 SQLite news.db")
		panic(err)
	}
	// defer Newsdb.Close()
	Newsdb.SetMaxOpenConns(1)
	if _, err := Newsdb.Exec(`
	CREATE TABLE IF NOT EXISTS news (
	id 				INTEGER PRIMARY KEY,
	zhihu 			TEXT CHECK(json_valid(zhihu)),
	coolapk			TEXT CHECK(json_valid(coolapk)),
	juejin			TEXT CHECK(json_valid(juejin)),
	weibo			TEXT CHECK(json_valid(weibo)),
	toutiao			TEXT CHECK(json_valid(toutiao)),
	thepaper		TEXT CHECK(json_valid(thepaper)),
	hackernews		TEXT CHECK(json_valid(hackernews)),
	clshot			TEXT CHECK(json_valid(clshot)),
	freebuf 		TEXT CHECK(json_valid(freebuf)),
	created_at 		TEXT DEFAULT (datetime('now','localtime'))
	);`); err != nil {
		log.Println("无法创建 news 表")
		panic(err)
	}
}

func _init_wordsclouddb() {
	var err error
	WordsClouddb, err = sql.Open("sqlite", filepath.Join(config.DataBaseDir, "/wordscloud.db"))
	if err != nil {
		log.Println("无法打开 SQLite database/wordscloud.db")
		panic(err)
	}
	// defer WordsClouddb.Close()
	WordsClouddb.SetMaxOpenConns(1)
	if _, err := WordsClouddb.Exec(`
	CREATE TABLE IF NOT EXISTS wordscloud (
	id 				INTEGER PRIMARY KEY,
	data 			TEXT CHECK(json_valid(data)),
	created_at 		TEXT DEFAULT (datetime('now','localtime'))
	);`); err != nil {
		log.Println("无法创建 wordscloud 表")
		panic(err)
	}
}

func _init_linkwords() {
	var err error
	LinkWords, err = sql.Open("sqlite", filepath.Join(config.DataBaseDir, "/linkwords.db"))
	if err != nil {
		log.Println("无法打开 SQLite database/linkwords.db")
		panic(err)
	}
	// defer LinkWords.Close()
	LinkWords.SetMaxOpenConns(1)
	if _, err := LinkWords.Exec(`
	CREATE TABLE IF NOT EXISTS linkwords (
	id 				INTEGER PRIMARY KEY,
	data 			TEXT CHECK(json_valid(data)),
	created_at 		TEXT DEFAULT (datetime('now','localtime'))
	);`); err != nil {
		log.Println("无法创建 linkwords 表")
		panic(err)
	}
}
