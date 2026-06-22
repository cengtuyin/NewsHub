package allstruct

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ModelInfo struct {
	Url   string
	Key   string
	Model string
}

type Pusher interface {
	GetInfo() *PushInfo
	GetSettings() map[string]any
	Push(message PushMessage) error
}

type PushInfo struct {
	Class string
	Url   string
}

type PushInfo_Webhook struct {
	PushInfo
	Header map[string]string
	Body   string
}

// GetInfo implements [Pusher].
func (push PushInfo_Webhook) GetInfo() *PushInfo {
	return &push.PushInfo
}

// GetSettings implements [Pusher].
func (push PushInfo_Webhook) GetSettings() map[string]any {
	return map[string]any{
		"Header": &push.Header,
		"Body":   &push.Body,
	}
}

// Push implements [Pusher].
func (push PushInfo_Webhook) Push(message PushMessage) error {
	var client http.Client
	req, _ := http.NewRequest("POST", message.Build(push.Url), bytes.NewBuffer([]byte(message.Build(push.Body))))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36 Edg/149.0.0.0")
	for k, v := range push.Header {
		req.Header.Set(k, v)
	}
	_, err := client.Do(req)
	return err
}

type PushInfo_Email struct {
	PushInfo
}

// GetInfo implements [Pusher].
func (push PushInfo_Email) GetInfo() *PushInfo {
	return &push.PushInfo
}

// GetSettings implements [Pusher].
func (push PushInfo_Email) GetSettings() map[string]any {
	return nil
}

// Push implements [Pusher].
func (push PushInfo_Email) Push(message PushMessage) error {
	return fmt.Errorf("尚未支持该推送方式")
}

type PushMessage struct {
	Type    string
	Title   string
	Content string
	Url     string
}

func (message PushMessage) Build(o string) string {
	o = strings.ReplaceAll(o, "${type}", url.QueryEscape(message.Type))
	o = strings.ReplaceAll(o, "${title}", url.QueryEscape(message.Title))
	o = strings.ReplaceAll(o, "${content}", url.QueryEscape(message.Content))
	o = strings.ReplaceAll(o, "${url}", url.QueryEscape(message.Url))
	return o
}

func Init() {

}
