package api

import (
	"encoding/xml"
	"fmt"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-wx-api/internal/config/wx"
	"github.com/xxl6097/go-wx-api/internal/ntfy"
	"io"
	"net/http"
)

func (this *Api) wxMessageRecv(openId string, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // 确保关闭 Body
	fmt.Printf("%s %v\n", string(body), err)
	var baseMsg struct {
		MsgType string `xml:"MsgType"`
	}
	if err = xml.Unmarshal(body, &baseMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch baseMsg.MsgType {
	case "text":
		var textMsg wx.TextMessage
		err = xml.Unmarshal(body, &textMsg)
		fmt.Println(textMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		break
	case "event":
		var eventMsg wx.EventMessage
		err = xml.Unmarshal(body, &eventMsg)
		fmt.Println(eventMsg)
		this.eventMessage(w, &eventMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		break
	case "image":
		break
	}
}

func (this *Api) eventMessage(w http.ResponseWriter, event *wx.EventMessage) {
	if event == nil {
		return
	}
	_ = ntfy.GetInstance().Publish(&ntfy.NtfyEventData{
		Title:   "sign",
		Topic:   "uclient",
		Message: event.EventKey,
	})
	switch event.EventKey {
	case "16:00:6f:83:35:e1":
		replyMsg := wx.CreateTextResponse(
			event.FromUserName,              // 接收方：发送消息的用户OpenID
			event.ToUserName,                // 发送方：公众号ID
			ntfy.GetInstance().GetMessage(), // 回复内容
		)
		// 2. 将结构体序列化为XML字节切片
		xmlData, err := xml.Marshal(replyMsg)
		if err != nil {
			glog.Printf("Error marshaling XML response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// 3. 设置响应头并返回XML
		w.Header().Set("Content-Type", "application/xml") // 务必设置为 application/xml[1,2](@ref)
		_, _ = w.Write(xmlData)
		break
	}
}

/*
{
    "button": [
        {
            "type": "click",
            "name": "打卡",
            "key": "16:00:6f:83:35:e1"
        },
        {
            "type": "click",
            "name": "查询",
            "key": "16:00:6f:83:35:e1"
        },
        {
            "name": "功能",
            "sub_button": [
                {
                    "type": "view",
                    "name": "504.7000",
                    "key": "WORK_003",
                    "url": "http://uuxia.cn:7000?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
                },
                {
                    "type": "view",
                    "name": "gz.7000",
                    "key": "WORK_004",
                    "url": "http://uuxia.cn:6633?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
                },
                {
                    "type": "view",
                    "name": "baidu",
                    "key": "WORK_005",
                    "url": "https://www.baidu.com"
                },
                {
                    "type": "click",
                    "name": "获取链接",
                    "key": "WORK_007"
                },
                {
                    "type": "view",
                    "name": "clife.7000",
                    "key": "WORK_002",
                    "url": "http://uuxia.cn:6615?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
                }
            ]
        }
    ]
}
*/
