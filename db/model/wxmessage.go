package model

import (
	"strconv"
	"strings"
)

type WXMessage struct {
	ToUserName   string       `json:"ToUserName" xml:"ToUserName"`
	FromUserName string       `json:"FromUserName" xml:"FromUserName"`
	CreateTime   int64        `json:"CreateTime" xml:"CreateTime"`
	MsgType      string       `json:"MsgType" xml:"MsgType"`
	Content      string       `json:"Content" xml:"Content"`
	Image        *WXImage     `json:"Image" xml:"Image"`
	Voice        *WXVoice     `json:"Voice" xml:"Voice"`
	Video        *WXVideo     `json:"Video" xml:"Video"`
	Music        *WXMusic     `json:"music" xml:"music"`
	Artical      []*WXArticle `json:"artical" xml:"artical"`
}

type WXImage struct {
	MediaId string `json:"MediaId" xml:"MediaId"`
}

type WXVoice struct {
	MediaId string `json:"MediaId" xml:"MediaId"`
}

type WXVideo struct {
	MediaId     string `json:"MediaId" xml:"MediaId"`
	Title       string `json:"Title" xml:"Title"`
	Description string `json:"Description" xml:"Description"`
}

type WXMusic struct {
	Title        string `json:"Title" xml:"Title"`
	Description  string `json:"Description" xml:"Description"`
	MusicUrl     string `json:"MusicUrl" xml:"MusicUrl"`
	HQMusicUrl   string `json:"HQMusicUrl" xml:"HQMusicUrl"`
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId"`
}

type WXArticle struct {
	Title       string `json:"Title" xml:"Title"`
	Description string `json:"Description" xml:"Description"`
	PicUrl      string `json:"PicUrl" xml:"PicUrl"`
	Url         string `json:"Url" xml:"Url"`
}

func (s *WXMessage) ToResponseXMLString() []byte {
	var ret strings.Builder
	ret.WriteString("<xml>\n")

	ret.WriteString("<ToUserName><![CDATA[")
	ret.WriteString(s.FromUserName)
	ret.WriteString("]]></ToUserName>\n")

	ret.WriteString("<FromUserName><![CDATA[")
	ret.WriteString(s.ToUserName)
	ret.WriteString("]]></FromUserName>\n")

	ret.WriteString("<CreateTime>")
	ret.WriteString(strconv.Itoa(int(s.CreateTime)))
	ret.WriteString("</CreateTime>\n")

	ret.WriteString("<MsgType><![CDATA[text]]></MsgType>\n")

	ret.WriteString("<Content><![CDATA[")
	ret.WriteString(s.Content)
	ret.WriteString("]]></Content>\n")

	ret.WriteString("</xml>")
	return []byte(ret.String())
}
