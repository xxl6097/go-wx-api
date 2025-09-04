package api

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/httpserver"
	"go-wx-api/internal/config"
	"go-wx-api/internal/config/wx"
	"go-wx-api/internal/ntfy"
	"go-wx-api/internal/u"
	"io"
	"log"
	"net/http"
	"time"
)

const yourToken = "het002402"

type Api struct {
	cfg             *config.Config
	cachedToken     string
	tokenExpireTime time.Time
}

func NewApi(cfg *config.Config) func(*mux.Router) {
	restApi := Api{cfg: cfg}
	token := restApi.GetToken()
	if token == "" {
		glog.Fatal("load app access token err:")
	} else {
		glog.Info("Token:", restApi.cachedToken)
	}
	return func(router *mux.Router) {
		wxRouter := router.NewRoute().Subrouter()
		wxRouter.HandleFunc("/api/wx/push", restApi.ApiWxPush)

		hRouter := router.NewRoute().Subrouter()
		httpserver.BasicAuth(hRouter, cfg.Username, cfg.Password)
		hRouter.HandleFunc("/api/hello", restApi.ApiHello)
	}
}

func (this *Api) GetToken() string {
	if this.cachedToken == "" || time.Now().After(this.tokenExpireTime) {
		// 缓存无效，重新获取
		tokenInfo, err := GetStableAccessToken(this.cfg.AppID, this.cfg.AppSecret, false)
		if err != nil {
			fmt.Printf("get stable access token err:%v\n", err)
			return ""
		}
		this.cachedToken = tokenInfo.AccessToken
		this.tokenExpireTime = time.Now().Add(time.Duration(tokenInfo.ExpiresIn-120) * time.Second) // 提前2分钟过期
	}
	return this.cachedToken
}

func wechatMessageHandler(w http.ResponseWriter, userMsg wx.EventMessage) {
	// 1. 构造回复消息
	replyMsg := wx.CreateTextResponse(
		userMsg.FromUserName, // 接收方：发送消息的用户OpenID
		userMsg.ToUserName,   // 发送方：公众号ID
		"您好，这是自动回复：\r\n"+
			"姓名：夏小力\r\n"+
			"性别：男\r\n"+
			"工作：码农", // 回复内容
	)
	// 2. 将结构体序列化为XML字节切片
	xmlData, err := xml.Marshal(replyMsg)
	if err != nil {
		log.Printf("Error marshaling XML response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// 3. 设置响应头并返回XML
	w.Header().Set("Content-Type", "application/xml") // 务必设置为 application/xml[1,2](@ref)
	w.Write(xmlData)
	log.Println("Text response sent successfully.")
}

func (this *Api) ApiWxPush(w http.ResponseWriter, r *http.Request) {
	glog.Printf("%s %s %s\n", r.Method, r.URL.String(), r.Proto)
	queryParams := r.URL.Query()
	echostr := queryParams.Get("echostr")
	nonce := queryParams.Get("nonce")
	openid := queryParams.Get("openid")
	signature := queryParams.Get("signature")
	timestamp := queryParams.Get("timestamp")
	switch r.Method {
	case http.MethodPost:
		ok := u.CheckSignature(signature, timestamp, nonce, yourToken)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			glog.Println("微信服务器验证失败: 签名无效")
			return
		}
		glog.Println("微信服务器验证成功")
		this.wxPushPost(openid, w, r)
		break
	case http.MethodGet:
		ok := u.CheckSignature(signature, timestamp, echostr, yourToken)
		if ok {
			// 若确认此次GET请求来自微信服务器，请原样返回 echostr 参数内容
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(echostr))
			glog.Println("微信服务器验证成功")
		} else {
			// 校验失败
			w.WriteHeader(http.StatusForbidden)
			glog.Println("微信服务器验证失败: 签名无效")
			return
		}
		break
	}
}
func (this *Api) wxPushPost(openId string, w http.ResponseWriter, r *http.Request) {
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
			event.FromUserName, // 接收方：发送消息的用户OpenID
			event.ToUserName,   // 发送方：公众号ID
			"您好，这是自动回复：\r\n"+
				"姓名：夏小力\r\n"+
				"性别：男\r\n"+
				"工作：码农", // 回复内容
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

func (this *Api) ApiHello(w http.ResponseWriter, r *http.Request) {
	glog.Printf("%s %s %s\n", r.Method, r.URL.String(), r.Proto)
	w.Write([]byte("hello world"))
}
