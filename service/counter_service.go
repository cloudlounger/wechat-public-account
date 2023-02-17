package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	openid := header.Get("x-wx-openid")
	if openid == "" {
		fmt.Println("-----------empty openid")
		w.WriteHeader(400)
		return
	}
	fmt.Println("-----------x-wx-openid", openid)
	msg := &model.WXMessage{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
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

	key := getKey(msg)
	if _, ok := cache.Load(key); !ok {
		cache.Store(key, struct{}{})
		pushQueue(msg)
	}
	quit := loopCheck(key, r.Context().Done())
	if quit {
		fmt.Println("-----------loopCheck quit, key:", key)
		return
	}
	v, ok := result.Load(key)
	if !ok {
		w.WriteHeader(500)
		return
	}
	fmt.Println("-----------Load:", v.(string))
	/*
		b, err = msg.ToResponseJsonStringWithOpenAI(v.(string))
		if err != nil {
			w.WriteHeader(500)
			return
		}
	*/
	msg.Content = v.(string)
	b = msg.ToResponseXMLString()
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	result.Delete(key)
	fmt.Println("-----------return", msg)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------hello")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
}

func SendAsync(msg *model.WXMessage) string {
	fmt.Println("begin SendMessage", msg.Content)
	word, err := defaultPayload.SendMessage(msg.Content)
	if err != nil {
		fmt.Println("[debug] defaultPayload.SendMessage failed, error", err)
		return ""
	}
	fmt.Println("finish SendMessage", word)
	return word
}

type WXCustomMessage struct {
	ToUser  string  `json:"touser"`
	Msgtype string  `json:"msgtype"`
	Text    *WXText `json:"text"`
}

type WXText struct {
	Content string `json:"content"`
}

func getKey(msg *model.WXMessage) string {
	return msg.FromUserName + strconv.Itoa(int(msg.CreateTime))
}
