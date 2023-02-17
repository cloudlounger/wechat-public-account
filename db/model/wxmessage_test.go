package model

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestWXMessage(t *testing.T) {
	msg := &WXMessage{
		ToUserName:   "11111",
		FromUserName: "22222",
		CreateTime:   87642323,
		MsgType:      "text",
		Content:      "i am ok",
	}
	b, err := xml.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	b = msg.ToResponseXMLString()
	fmt.Println(string(b))
}
