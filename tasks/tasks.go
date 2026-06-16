package tasks

import (
	"log"
	"net/http"
	"newshub/config"
	"time"
)

func Init() {
	go _init()
}

func _init() {
	for !config.Ok {
		time.Sleep(time.Duration(1) * time.Second)
	}

	// 填充缓存
	var client http.Client
	client.Get("http://localhost:51516/api/newsnow2")
	client.Get("http://localhost:51516/api/wordscloud?length=7")
	client.Get("http://localhost:51516/api/wordsin?length=2")

	log.Println("数据就绪，计划任务正在挑选一个吉利的时机入场")

	time.Sleep(time.Duration(60-time.Now().Second()) * time.Second)
	go task_30minutes()
	go task_72hours()
	// go task_nextday()

	log.Println("已完整载入")
}

func task_30minutes() {
	var client http.Client
	for true {
		time.Sleep(time.Duration(30) * time.Minute)
		log.Println("每 30 分钟次计划任务激活")
		client.Get("http://localhost:51516/api/newsnow2")
		client.Get("http://localhost:51516/api/wordscloud?length=7")
		client.Get("http://localhost:51516/api/wordsin?length=2")
	}
}

func task_72hours() {
	for true {
		time.Sleep(time.Duration(72) * time.Hour)
		log.Println("每 72 小时次计划任务激活")
		config.Logined = []string{}
	}
}

func task_nextday() {
	for true {
		time.Sleep(time.Duration(23-time.Now().Hour()) * time.Hour)
		time.Sleep(time.Duration(59-time.Now().Minute()) * time.Minute)
		time.Sleep(time.Duration(59-time.Now().Second()) * time.Second)
		log.Println("每 00:00 计划任务激活")
	}
}
