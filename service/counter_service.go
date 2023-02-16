package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// CounterHandler 计数器接口
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		counter, err := getCurrentCounter()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = counter.Count
		}
	} else if r.Method == http.MethodPost {
		count, err := modifyCounter(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = count
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	action, err := getAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCounter(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = clearCounter()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
		createdAt = time.Now()
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.Imp.UpsertCounter(counter)
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func clearCounter() error {
	return dao.Imp.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	counter, err := dao.Imp.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	action, ok := body["action"]
	if !ok {
		return "", fmt.Errorf("缺少 action 参数")
	}

	return action.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WXMessageHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	body := r.Body
	defer body.Close()
	//appid := header.Get("x-wx-from-appid")
	//if appid == "" {
	//	fmt.Println("-----------empty appid")
	//	w.WriteHeader(400)
	//	return
	//}
	openid := header.Get("x-wx-openid")
	if openid == "" {
		fmt.Println("-----------empty openid")
		w.WriteHeader(400)
		return
	}
	fmt.Println("-----------x-wx-openid", openid)
	msg := &model.WXMessage{}
	b, err := ioutil.ReadAll(body)
	if err != nil && err != io.EOF {
		fmt.Println("-----------ReadAll failed", err)
		w.WriteHeader(400)
		return
	}
	err = json.Unmarshal(b, msg)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	fmt.Printf("-----------success %+v\n", msg)
	/*
		xmlb := msg.ToResponseXMLString()
		fmt.Println("-----------xml", xmlb)
		_, err = w.Write(xmlb)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	*/
	//b, err = msg.ToResponseJsonString()
	//if err != nil {
	//	w.WriteHeader(500)
	//	return
	//}
	go func() {
		SendAsync(msg)
	}()
	word := "异步调用openai中, 请耐心等待"
	b, err = msg.ToResponseJsonStringWithOpenAI(word)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(b)
	w.Header().Set("content-type", "application/json")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------hello")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
}

func SendAsync(msg *model.WXMessage) {
	word, err := defaultPayload.SendMessage(msg.Content)
	if err != nil {
		fmt.Println("[debug] defaultPayload.SendMessage failed, error", err)
		return
	}
	customMsg := &WXCustomMessage{
		ToUser:  msg.FromUserName,
		Msgtype: "text",
		Text:    WXText{Content: word},
	}
	payloadBytes, err := json.Marshal(customMsg)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", "http://api.weixin.qq.com/cgi-bin/message/custom/send", body)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[debug] error", err)
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("[debug] error status code, code", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
}

type WXCustomMessage struct {
	ToUser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Text    WXText `json:"text"`
}

type WXText struct {
	Content string `json:"content"`
}
